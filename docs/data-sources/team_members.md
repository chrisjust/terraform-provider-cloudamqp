---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_members"
description: |-
  Retrieve all team members of your CloudAMQP account.
---

# cloudamqp_team_members

This data source allows you to retrieve all team members on your CloudAMQP account.

## Example Usage

```hcl
data "cloudamqp_team_members" "all" {}

output "member_emails" {
  value = [for member in data.cloudamqp_team_members.all.members : member.email]
}
```

## Argument Reference

No arguments are supported for this data source.

## Attributes Reference

All attributes references are computed

* `id`      - The identifier for this resource, always set to `team`.
* `members` - An array of team members. Each `members` block consists of the field documented below.

___

The `members` block consists of:

* `user_id`          - User identifier in UUID format.
* `email`            - Email address of the team member.
* `tfa_auth_enabled` - Whether two factor authentication is enabled for the team member.
* `roles`            - Roles of the team member.

## Import

`cloudamqp_team_members` can be imported using a fixed identifier, shown below.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = data.cloudamqp_team_members.all
  id = "team"
}
```

Or use Terraform CLI:

`terraform import data.cloudamqp_team_members.all team`

[CloudAMQP API list team members]: https://docs.cloudamqp.com/#tag/team/GET/team