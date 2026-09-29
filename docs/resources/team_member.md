---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_member"
description: |-
  Invite and manage a team member on your CloudAMQP account.
---

# cloudamqp_team_member

This resource allows you to invite and manage a team member on your CloudAMQP account. The member
only appears in the account team after the invitation has been accepted; until then the user id is
empty and role and tag changes are deferred until after acceptance.

## Example Usage

```hcl
resource "cloudamqp_team_member" "member" {
  email = "member@example.com"
  role  = "devops"
  tags  = ["oncall", "platform"]
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) Email address of the team member. Changing this forces a new resource to be created.
* `role`  - (Optional) Role of the team member. Valid options are: `admin`, `devops`, `member`,
            `monitor`, `billing manager`. Default set to `member`.
* `tags`  - (Optional) Tags for the team member. Changing this while the invitation is still
            pending is recorded and applied together after the invitee accepts.

## Attributes Reference

All attributes references are computed

* `id` - User identifier in UUID format. Empty until the invitee accepts the invitation.

## Import

`cloudamqp_team_member` can be imported using the email address of the team member. To retrieve the email
address, use the [CloudAMQP API list team members].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_team_member.member
  id = "member@example.com"
}
```

Or use Terraform CLI:

`terraform import cloudamqp_team_member.member member@example.com`

[CloudAMQP API list team members]: https://docs.cloudamqp.com/#tag/team/GET/team