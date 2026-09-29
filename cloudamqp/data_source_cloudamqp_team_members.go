package cloudamqp

import (
	"context"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeamMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamMembersRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of team members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User identifier in UUID format",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address of the team member",
						},
						"tfa_auth_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether two factor authentication is enabled for the team member",
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Roles of the team member",
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Tags of the team member",
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamMembersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*api.API)
	data, err := client.ListTeamMembers(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("team")

	members := make([]map[string]any, len(data))
	for k, v := range data {
		members[k] = readTeamMember(v)
	}

	if err = d.Set("members", members); err != nil {
		return diag.Errorf("error setting members for resource %s: %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func readTeamMember(data model.TeamMemberResponse) map[string]any {
	tags := data.Tags
	if tags == nil {
		tags = []string{}
	}
	return map[string]any{
		"user_id":          data.ID,
		"email":            data.Email,
		"tfa_auth_enabled": data.TfaAuthEnabled,
		"roles":            data.Roles,
		"tags":             tags,
	}
}
