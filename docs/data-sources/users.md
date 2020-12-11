---
page_title: "octopusdeploy_users Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing users.
---

# Data Source `octopusdeploy_users`

Provides information about existing users.

## Example Usage

```terraform
data "octopusdeploy_users" "example" {
  ids  = ["Users-123", "Users-321"]
  skip = 5
  take = 100
}
```

## Schema

### Optional

- **filter** (String, Optional) A filter with which to search.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **user** (Block List) A list of users that match the filter(s). (see [below for nested schema](#nestedblock--user))

<a id="nestedblock--user"></a>
### Nested Schema for `user`

Read-only:

- **can_password_be_edited** (Boolean, Read-only)
- **display_name** (String, Read-only) The display name of this resource.
- **email_address** (String, Read-only) The email address of this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **identity** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--user--identity))
- **is_active** (Boolean, Read-only)
- **is_requestor** (Boolean, Read-only)
- **is_service** (Boolean, Read-only)
- **password** (String, Read-only) The password associated with this resource.
- **username** (String, Read-only) The username associated with this resource.

<a id="nestedatt--user--identity"></a>
### Nested Schema for `user.identity`

- **claim** (Set of Object) (see [below for nested schema](#nestedobjatt--user--identity--claim))
- **provider** (String)

<a id="nestedobjatt--user--identity--claim"></a>
### Nested Schema for `user.identity.claim`

- **is_identifying_claim** (Boolean)
- **name** (String)
- **value** (String)


