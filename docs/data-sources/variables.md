---
page_title: "octopusdeploy_variables Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing variables.
---

# Data Source `octopusdeploy_variables`

Provides information about existing variables.

## Example Usage

```terraform
data "octopusdeploy_variables" "example" {
}
```

## Schema

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **variable** (Block List) A list of variables that match the filter(s). (see [below for nested schema](#nestedblock--variable))

<a id="nestedblock--variable"></a>
### Nested Schema for `variable`

Read-only:

- **description** (String, Read-only) The description of this resource.
- **encrypted_value** (String, Read-only)
- **is_sensitive** (Boolean, Read-only) Indicates whether or not this resource is considered sensitive and should be kept secret.
- **key_fingerprint** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **pgp_key** (String, Read-only)
- **project_id** (String, Read-only)
- **prompt** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--variable--prompt))
- **scope** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--variable--scope))
- **sensitive_value** (String, Read-only)
- **type** (String, Read-only) The type of variable represented by this resource. Valid types are `AmazonWebServicesAccount`, `AzureAccount`, `Certificate`, `Sensitive`, `String`, or `WorkerPool`.
- **value** (String, Read-only)

<a id="nestedatt--variable--prompt"></a>
### Nested Schema for `variable.prompt`

- **description** (String)
- **is_required** (Boolean)
- **label** (String)


<a id="nestedatt--variable--scope"></a>
### Nested Schema for `variable.scope`

- **actions** (List of String)
- **channels** (List of String)
- **environments** (List of String)
- **machines** (List of String)
- **roles** (List of String)
- **tenant_tags** (List of String)


