---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_member"
description: |-
  Invite and manage a team member on your CloudAMQP account.
---

# cloudamqp_team_member

This resource allows you to invite and manage a team member on your CloudAMQP account.

## Example Usage

```hcl
resource "cloudamqp_team_member" "member" {
  email = "member@example.com"
  role  = "devops"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) Email address of the team member. Changing this forces a new resource to be created.
* `role`  - (Optional) Role of the team member. Valid options are: `admin`, `devops`, `member`,
            `monitor`, `billing manager`. Default set to `member`.

## Attributes Reference

All attributes references are computed

* `id`  - The identifier for this resource, set to the email address of the team member.

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