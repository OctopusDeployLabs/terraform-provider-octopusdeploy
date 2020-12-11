---
page_title: "octopusdeploy_user Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages users in Octopus Deploy.
---

# Resource `octopusdeploy_user`

This resource manages users in Octopus Deploy.



## Schema

### Required

- **display_name** (String, Required) The display name of this resource.
- **username** (String, Required) The username associated with this resource.

### Optional

- **email_address** (String, Optional) The email address of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **identity** (Block Set) (see [below for nested schema](#nestedblock--identity))
- **is_active** (Boolean, Optional)
- **is_service** (Boolean, Optional)
- **password** (String, Optional) The password associated with this resource.

### Read-only

- **can_password_be_edited** (Boolean, Read-only)
- **is_requestor** (Boolean, Read-only)

<a id="nestedblock--identity"></a>
### Nested Schema for `identity`

Optional:

- **claim** (Block Set) (see [below for nested schema](#nestedblock--identity--claim))
- **provider** (String, Optional)

<a id="nestedblock--identity--claim"></a>
### Nested Schema for `identity.claim`

Required:

- **is_identifying_claim** (Boolean, Required)
- **name** (String, Required) The name of this resource.
- **value** (String, Required)


