variable "team_member_email" {
  description = "Email address of the team member to invite"
  type        = string
}

variable "team_member_role" {
  description = "Role of the team member (admin, devops, member, monitor, billing manager)"
  type        = string
  default     = "member"
}