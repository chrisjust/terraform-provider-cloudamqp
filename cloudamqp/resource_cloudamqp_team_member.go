package cloudamqp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &teamMemberResource{}
	_ resource.ResourceWithConfigure = &teamMemberResource{}
)

func NewTeamMemberResource() resource.Resource {
	return &teamMemberResource{}
}

type teamMemberResource struct {
	client *api.API
}

type teamMemberResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
	Tags  types.Set    `tfsdk:"tags"`
}

func (r *teamMemberResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*api.API)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *teamMemberResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "cloudamqp_team_member"
}

func (r *teamMemberResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Invite and manage a team member on your CloudAMQP account. The member only " +
			"appears in the account team after the invitation is accepted; until then the user id " +
			"is empty and role changes are deferred.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "User identifier in UUID format. Empty until the invitee accepts the invitation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email address of the team member",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Tags for the team member",
			},
			"role": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Role of the team member. Valid options are: admin, devops, member, monitor, billing manager",
				Validators: []validator.String{
					stringvalidator.OneOf("admin", "devops", "member", "monitor", "billing manager"),
				},
			},
		},
	}
}

func (r *teamMemberResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan teamMemberResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	email := plan.Email.ValueString()
	role := plan.Role.ValueString()
	if role == "" {
		role = "member"
	}

	params := model.InviteTeamMemberRequest{
		Email: email,
		Role:  role,
		Tags:  setToStringSlice(plan.Tags),
	}
	if _, err := r.client.InviteTeamMember(ctx, params); err != nil {
		if isAlreadyInvited(err) {
			tflog.Info(ctx, fmt.Sprintf("team member %s already invited, treating as converged", email))
		} else {
			response.Diagnostics.AddError(
				"Failed to invite team member",
				fmt.Sprintf("Could not invite team member %s: %s", email, err),
			)
			return
		}
	}

	plan.Id = types.StringNull()
	plan.Role = types.StringValue(role)
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *teamMemberResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state teamMemberResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	member, err := r.findMember(ctx, state.Email.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to list team members",
			fmt.Sprintf("Could not list team members: %s", err),
		)
		return
	}

	if member != nil {
		state.Id = types.StringValue(member.ID)
		state.Role = types.StringValue(primaryRole(member.Roles))
		state.Tags = stringSliceToSet(member.Tags)
		response.Diagnostics.Append(response.State.Set(ctx, &state)...)
		return
	}

	if state.Id.IsNull() || state.Id.ValueString() == "" {
		tflog.Info(ctx, fmt.Sprintf("team member %s invite still pending, keeping state unchanged", state.Email.ValueString()))
		response.Diagnostics.Append(response.State.Set(ctx, &state)...)
		return
	}

	response.State.RemoveResource(ctx)
}

func (r *teamMemberResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan teamMemberResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.Id.IsNull() || plan.Id.ValueString() == "" {
		tflog.Info(ctx, fmt.Sprintf("team member %s invite still pending, deferring role update", plan.Email.ValueString()))
		response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
		return
	}

	params := model.UpdateTeamMemberRequest{
		Role: plan.Role.ValueString(),
		Tags: setToStringSlice(plan.Tags),
	}
	if _, err := r.client.UpdateTeamMember(ctx, plan.Id.ValueString(), params); err != nil {
		response.Diagnostics.AddError(
			"Failed to update team member",
			fmt.Sprintf("Could not update team member %s: %s", plan.Email.ValueString(), err),
		)
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *teamMemberResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state teamMemberResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	email := state.Email.ValueString()
	if _, err := r.client.RemoveTeamMember(ctx, email); err != nil {
		member, listErr := r.findMember(ctx, email)
		if listErr != nil || member != nil {
			response.Diagnostics.AddError(
				"Failed to remove team member",
				fmt.Sprintf("Could not remove team member %s: %s", email, err),
			)
			return
		}
		tflog.Info(ctx, fmt.Sprintf("team member %s not found, treating as removed", email))
	}
}

func (r *teamMemberResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	response.State.SetAttribute(ctx, path.Root("email"), request.ID)
	response.State.SetAttribute(ctx, path.Root("id"), types.StringNull())
	response.State.SetAttribute(ctx, path.Root("role"), types.StringUnknown())
	response.State.SetAttribute(ctx, path.Root("tags"), types.SetNull(types.StringType))
}

func (r *teamMemberResource) findMember(ctx context.Context, email string) (*model.TeamMemberResponse, error) {
	members, err := r.client.ListTeamMembers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range members {
		if strings.EqualFold(members[i].Email, strings.TrimSpace(email)) {
			return &members[i], nil
		}
	}
	return nil, nil
}

func isAlreadyInvited(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already invited")
}

func primaryRole(roles []string) string {
	if len(roles) == 0 {
		return "member"
	}
	return roles[0]
}

func setToStringSlice(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	elements := set.Elements()
	tags := make([]string, 0, len(elements))
	for _, element := range elements {
		if tag, ok := element.(types.String); ok {
			tags = append(tags, tag.ValueString())
		}
	}
	return tags
}

func stringSliceToSet(tags []string) types.Set {
	if len(tags) == 0 {
		return types.SetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		elements = append(elements, types.StringValue(tag))
	}
	value, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}
	return value
}
