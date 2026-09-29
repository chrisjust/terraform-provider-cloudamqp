package team

type TeamMemberResponse struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	TfaAuthEnabled bool     `json:"tfa_auth_enabled"`
	Roles          []string `json:"roles"`
	Tags           []string `json:"tags"`
}

type InviteTeamMemberRequest struct {
	Email string   `url:"email"`
	Role  string   `url:"role,omitempty"`
	Tags  []string `url:"tags" del:","`
}

type UpdateTeamMemberRequest struct {
	Role string   `url:"role,omitempty"`
	Tags []string `url:"tags" del:","`
}

type RemoveTeamMemberRequest struct {
	Email string `url:"email"`
}
