---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_lifecycles Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing lifecycles.
---

# octopusdeploy_lifecycles (Data Source)

Provides information about existing lifecycles.

## Example Usage

```terraform
data "octopusdeploy_lifecycles" "example" {
  ids          = ["Lifecycles-123", "Lifecycles-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ids` (List of String) A list of lifecycle IDs to filter by.
- `partial_name` (String) A partial name to filter lifecycles by.
- `skip` (Number) A filter to specify the number of items to skip in the response.
- `space_id` (String) The space ID associated with this lifecycle.
- `take` (Number) A filter to specify the number of items to take (or return) in the response.

### Read-Only

- `id` (String) The ID of the lifecycle.
- `lifecycles` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles))

<a id="nestedatt--lifecycles"></a>
### Nested Schema for `lifecycles`

Read-Only:

- `description` (String) The description of the lifecycle.
- `id` (String) The ID of the lifecycle.
- `name` (String) The name of the lifecycle.
- `phase` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles--phase))
- `release_retention_policy` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles--release_retention_policy))
- `space_id` (String) The space ID associated with this lifecycle.
- `tentacle_retention_policy` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles--tentacle_retention_policy))

<a id="nestedatt--lifecycles--phase"></a>
### Nested Schema for `lifecycles.phase`

Read-Only:

- `automatic_deployment_targets` (List of String) The automatic deployment targets for this phase.
- `id` (String) The ID of the phase.
- `is_optional_phase` (Boolean) Whether this phase is optional.
- `is_priority_phase` (Boolean) Deployments will be prioritized in this phase
- `minimum_environments_before_promotion` (Number) The minimum number of environments before promotion.
- `name` (String) The name of the phase.
- `optional_deployment_targets` (List of String) The optional deployment targets for this phase.
- `release_retention_policy` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles--phase--release_retention_policy))
- `tentacle_retention_policy` (Attributes List) (see [below for nested schema](#nestedatt--lifecycles--phase--tentacle_retention_policy))

<a id="nestedatt--lifecycles--phase--release_retention_policy"></a>
### Nested Schema for `lifecycles.phase.release_retention_policy`

Read-Only:

- `quantity_to_keep` (Number) The quantity of releases to keep.
- `should_keep_forever` (Boolean) Whether releases should be kept forever.
- `unit` (String) The unit of time for the retention policy.


<a id="nestedatt--lifecycles--phase--tentacle_retention_policy"></a>
### Nested Schema for `lifecycles.phase.tentacle_retention_policy`

Read-Only:

- `quantity_to_keep` (Number) The quantity of releases to keep.
- `should_keep_forever` (Boolean) Whether releases should be kept forever.
- `unit` (String) The unit of time for the retention policy.



<a id="nestedatt--lifecycles--release_retention_policy"></a>
### Nested Schema for `lifecycles.release_retention_policy`

Read-Only:

- `quantity_to_keep` (Number) The quantity of releases to keep.
- `should_keep_forever` (Boolean) Whether releases should be kept forever.
- `unit` (String) The unit of time for the retention policy.


<a id="nestedatt--lifecycles--tentacle_retention_policy"></a>
### Nested Schema for `lifecycles.tentacle_retention_policy`

Read-Only:

- `quantity_to_keep` (Number) The quantity of releases to keep.
- `should_keep_forever` (Boolean) Whether releases should be kept forever.
- `unit` (String) The unit of time for the retention policy.


