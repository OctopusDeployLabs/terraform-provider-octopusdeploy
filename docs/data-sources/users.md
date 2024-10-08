---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_users Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing users.
---

# octopusdeploy_users (Data Source)

Provides information about existing users.

## Example Usage

```terraform
data "octopusdeploy_users" "example" {
  ids  = ["Users-123", "Users-321"]
  skip = 5
  take = 100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (String) A filter search by username, display name or email
- `ids` (List of String) A filter to search by a list of IDs.
- `skip` (Number) A filter to specify the number of items to skip in the response.
- `space_id` (String, Deprecated) The space ID associated with this user.
- `take` (Number) A filter to specify the number of items to take (or return) in the response.

### Read-Only

- `id` (String) The unique ID for this resource.
- `users` (Attributes List) (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Required:

- `display_name` (String) The display name of this resource.
- `username` (String) The username associated with this resource.

Optional:

- `can_password_be_edited` (Boolean) Specifies whether or not the password can be edited.
- `email_address` (String) The email address of this resource.
- `identity` (Attributes Set) The identities associated with the user. (see [below for nested schema](#nestedatt--users--identity))
- `is_active` (Boolean) Specifies whether or not the user is active.
- `is_requestor` (Boolean) Specifies whether or not the user is the requestor.
- `is_service` (Boolean) Specifies whether or not the user is a service account.

Read-Only:

- `id` (String) The unique ID for this resource.

<a id="nestedatt--users--identity"></a>
### Nested Schema for `users.identity`

Read-Only:

- `claim` (Attributes Set) The claim associated with the identity. (see [below for nested schema](#nestedatt--users--identity--claim))
- `provider` (String) The identity provider.

<a id="nestedatt--users--identity--claim"></a>
### Nested Schema for `users.identity.claim`

Required:

- `name` (String) The name of this resource.
- `value` (String) The value of this resource.

Optional:

- `is_identifying_claim` (Boolean) Specifies whether or not the claim is an identifying claim.


