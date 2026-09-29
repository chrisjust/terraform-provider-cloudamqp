package team

type TeamMemberResponse struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	TfaAuthEnabled bool     `json:"tfa_auth_enabled"`
	Roles          []string `json:"roles"`
}

type InviteTeamMemberRequest struct {
	Email string `url:"email"`
	Role  string `url:"role,omitempty"`
	Tags  string `url:"tags,omitempty"`
}

type UpdateTeamMemberRequest struct {
	Role string `url:"role,omitempty"`
	Tags string `url:"tags,omitempty"`
}
