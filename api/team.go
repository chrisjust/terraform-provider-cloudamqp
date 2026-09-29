package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ListTeamMembers(ctx context.Context) ([]model.TeamMemberResponse, error) {
	var (
		data   []model.TeamMemberResponse
		failed map[string]any
		path   = "/api/team"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ListTeamMembers",
		resourceName: "Team",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) InviteTeamMember(ctx context.Context, params model.InviteTeamMemberRequest) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = "/api/team/invite"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyForm(params), retryRequest{
		functionName: "InviteTeamMember",
		resourceName: "Team",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) RemoveTeamMember(ctx context.Context, email string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = "/api/team/remove"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s email=%s", path, email))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyForm(model.RemoveTeamMemberRequest{Email: email}), retryRequest{
		functionName: "RemoveTeamMember",
		resourceName: "Team",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) UpdateTeamMember(ctx context.Context, userID string, params model.UpdateTeamMemberRequest) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/team/%s", userID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyForm(params), retryRequest{
		functionName: "UpdateTeamMember",
		resourceName: "Team",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
