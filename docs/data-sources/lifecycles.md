---
page_title: "octopusdeploy_lifecycles Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing lifecycles.
---

# Data Source `octopusdeploy_lifecycles`

Provides information about existing lifecycles.



## Schema

### Read-only

- **description** (String, Read-only) The description of this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **phase** (List of Object, Read-only) (see [below for nested schema](#nestedatt--phase))
- **release_retention_policy** (List of Object, Read-only) (see [below for nested schema](#nestedatt--release_retention_policy))
- **space_id** (String, Read-only) The space ID associated with this resource.
- **tentacle_retention_policy** (List of Object, Read-only) (see [below for nested schema](#nestedatt--tentacle_retention_policy))

<a id="nestedatt--phase"></a>
### Nested Schema for `phase`

- **automatic_deployment_targets** (List of String)
- **id** (String)
- **is_optional_phase** (Boolean)
- **minimum_environments_before_promotion** (Number)
- **name** (String)
- **optional_deployment_targets** (List of String)
- **release_retention_policy** (List of Object) (see [below for nested schema](#nestedobjatt--phase--release_retention_policy))
- **tentacle_retention_policy** (List of Object) (see [below for nested schema](#nestedobjatt--phase--tentacle_retention_policy))

<a id="nestedobjatt--phase--release_retention_policy"></a>
### Nested Schema for `phase.release_retention_policy`

- **quantity_to_keep** (Number)
- **should_keep_forever** (Boolean)
- **unit** (String)


<a id="nestedobjatt--phase--tentacle_retention_policy"></a>
### Nested Schema for `phase.tentacle_retention_policy`

- **quantity_to_keep** (Number)
- **should_keep_forever** (Boolean)
- **unit** (String)



<a id="nestedatt--release_retention_policy"></a>
### Nested Schema for `release_retention_policy`

- **quantity_to_keep** (Number)
- **should_keep_forever** (Boolean)
- **unit** (String)


<a id="nestedatt--tentacle_retention_policy"></a>
### Nested Schema for `tentacle_retention_policy`

- **quantity_to_keep** (Number)
- **should_keep_forever** (Boolean)
- **unit** (String)


