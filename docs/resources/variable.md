---
page_title: "octopusdeploy_variable Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_variable`





## Schema

### Required

- **name** (String, Required) The name of this resource.
- **project_id** (String, Required)
- **type** (String, Required) The type of variable represented by this resource. Valid types are `AmazonWebServicesAccount`, `AzureAccount`, `Certificate`, `Sensitive`, `String`, or `WorkerPool`.

### Optional

- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The ID of this resource.
- **is_sensitive** (Boolean, Optional) Indicates whether or not this resource is considered sensitive and should be kept secret.
- **pgp_key** (String, Optional)
- **prompt** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--prompt))
- **scope** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--scope))
- **sensitive_value** (String, Optional)
- **value** (String, Optional)

### Read-only

- **encrypted_value** (String, Read-only)
- **key_fingerprint** (String, Read-only)

<a id="nestedblock--prompt"></a>
### Nested Schema for `prompt`

Optional:

- **description** (String, Optional) The description of this resource.
- **is_required** (Boolean, Optional)
- **label** (String, Optional)


<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Optional:

- **actions** (List of String, Optional) The scope of the variable value.
- **channels** (List of String, Optional) The scope of the variable value.
- **environments** (List of String, Optional) The scope of the variable value.
- **machines** (List of String, Optional) The scope of the variable value.
- **roles** (List of String, Optional) The scope of the variable value.
- **tenant_tags** (List of String, Optional) The scope of the variable value.


