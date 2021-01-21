---
page_title: "octopusdeploy_projects Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing projects.
---

# Data Source `octopusdeploy_projects`

Provides information about existing projects.

## Example Usage

```terraform
data "octopusdeploy_projects" "example" {
  cloned_from_project_id = "Projects-456"
  ids                    = ["Projects-123", "Projects-321"]
  is_clone               = true
  name                   = "Default"
  partial_name           = "Defau"
  skip                   = 5
  take                   = 100
}
```

## Schema

### Optional

- **cloned_from_project_id** (String, Optional) A filter to search for cloned resources by a project ID.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **is_clone** (Boolean, Optional) A filter to search for cloned resources.
- **name** (String, Optional) A filter to search by name.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.
- **projects** (Block List) A list of projects that match the filter(s). (see [below for nested schema](#nestedblock--projects))

<a id="nestedblock--projects"></a>
### Nested Schema for `projects`

Read-only:

- **allow_deployments_to_no_targets** (Boolean, Read-only, Deprecated)
- **auto_create_release** (Boolean, Read-only)
- **auto_deploy_release_overrides** (List of String, Read-only)
- **cloned_from_project_id** (String, Read-only)
- **connectivity_policy** (List of Object, Read-only) (see [below for nested schema](#nestedatt--projects--connectivity_policy))
- **default_guided_failure_mode** (String, Read-only)
- **default_to_skip_if_already_installed** (Boolean, Read-only)
- **deployment_changes_template** (String, Read-only)
- **deployment_process_id** (String, Read-only)
- **description** (String, Read-only) The description of this resource.
- **discrete_channel_release** (Boolean, Read-only) Treats releases of different channels to the same environment as a separate deployment dimension
- **extension_settings** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--projects--extension_settings))
- **id** (String, Read-only) The unique ID for this resource.
- **included_library_variable_sets** (List of String, Read-only)
- **is_disabled** (Boolean, Read-only)
- **is_discrete_channel_release** (Boolean, Read-only) Treats releases of different channels to the same environment as a separate deployment dimension
- **is_version_controlled** (Boolean, Read-only)
- **lifecycle_id** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **project_group_id** (String, Read-only)
- **release_creation_strategy** (List of Object, Read-only) (see [below for nested schema](#nestedatt--projects--release_creation_strategy))
- **release_notes_template** (String, Read-only)
- **slug** (String, Read-only)
- **space_id** (String, Read-only) The space ID associated with this resource.
- **templates** (List of Object, Read-only) (see [below for nested schema](#nestedatt--projects--templates))
- **tenanted_deployment_participation** (String, Read-only) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **variable_set_id** (String, Read-only)
- **version_control_settings** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--projects--version_control_settings))
- **versioning_strategy** (Set of Object, Read-only) (see [below for nested schema](#nestedatt--projects--versioning_strategy))

<a id="nestedatt--projects--connectivity_policy"></a>
### Nested Schema for `projects.connectivity_policy`

- **allow_deployments_to_no_targets** (Boolean)
- **exclude_unhealthy_targets** (Boolean)
- **skip_machine_behavior** (String)
- **target_roles** (List of String)


<a id="nestedatt--projects--extension_settings"></a>
### Nested Schema for `projects.extension_settings`

- **extension_id** (String)
- **values** (List of String)


<a id="nestedatt--projects--release_creation_strategy"></a>
### Nested Schema for `projects.release_creation_strategy`

- **channel_id** (String)
- **release_creation_package** (List of Object) (see [below for nested schema](#nestedobjatt--projects--release_creation_strategy--release_creation_package))
- **release_creation_package_step_id** (String)

<a id="nestedobjatt--projects--release_creation_strategy--release_creation_package"></a>
### Nested Schema for `projects.release_creation_strategy.release_creation_package`

- **deployment_action** (String)
- **package_reference** (String)



<a id="nestedatt--projects--templates"></a>
### Nested Schema for `projects.templates`

- **default_value** (String)
- **display_settings** (Map of String)
- **help_text** (String)
- **id** (String)
- **label** (String)
- **name** (String)


<a id="nestedatt--projects--version_control_settings"></a>
### Nested Schema for `projects.version_control_settings`

- **default_branch** (String)
- **password** (String)
- **url** (String)
- **username** (String)


<a id="nestedatt--projects--versioning_strategy"></a>
### Nested Schema for `projects.versioning_strategy`

- **donor_package** (List of Object) (see [below for nested schema](#nestedobjatt--projects--versioning_strategy--donor_package))
- **donor_package_step_id** (String)
- **template** (String)

<a id="nestedobjatt--projects--versioning_strategy--donor_package"></a>
### Nested Schema for `projects.versioning_strategy.donor_package`

- **deployment_action** (String)
- **package_reference** (String)


