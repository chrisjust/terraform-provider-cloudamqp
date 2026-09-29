resource "cloudamqp_team_member" "member" {
  email = var.team_member_email
  role  = var.team_member_role
}