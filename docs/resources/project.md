---
page_title: "octopusdeploy_project Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages projects in Octopus Deploy.
---

# Resource `octopusdeploy_project`

This resource manages projects in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **allow_deployments_to_no_targets** (Boolean, Optional, Deprecated)
- **auto_create_release** (Boolean, Optional)
- **auto_deploy_release_overrides** (List of String, Optional)
- **cloned_from_project_id** (String, Optional)
- **connectivity_policy** (Block List) (see [below for nested schema](#nestedblock--connectivity_policy))
- **default_guided_failure_mode** (String, Optional)
- **default_to_skip_if_already_installed** (Boolean, Optional)
- **deployment_changes_template** (String, Optional)
- **description** (String, Optional) The description of this resource.
- **discrete_channel_release** (Boolean, Optional) Treats releases of different channels to the same environment as a separate deployment dimension
- **id** (String, Optional) The unique ID for this resource.
- **included_library_variable_sets** (List of String, Optional)
- **is_disabled** (Boolean, Optional)
- **is_discrete_channel_release** (Boolean, Optional) Treats releases of different channels to the same environment as a separate deployment dimension
- **is_version_controlled** (Boolean, Optional)
- **lifecycle_id** (String, Optional)
- **project_group_id** (String, Optional)
- **release_creation_strategy** (Block List, Max: 1) (see [below for nested schema](#nestedblock--release_creation_strategy))
- **release_notes_template** (String, Optional)
- **space_id** (String, Optional) The space ID associated with this resource.
- **templates** (Block List) (see [below for nested schema](#nestedblock--templates))
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **version_control_settings** (Block Set) (see [below for nested schema](#nestedblock--version_control_settings))

### Read-only

- **deployment_process_id** (String, Read-only)
- **extension_settings** (Block Set) (see [below for nested schema](#nestedblock--extension_settings))
- **slug** (String, Read-only)
- **variable_set_id** (String, Read-only)
- **versioning_strategy** (Block Set) (see [below for nested schema](#nestedblock--versioning_strategy))

<a id="nestedblock--connectivity_policy"></a>
### Nested Schema for `connectivity_policy`

Optional:

- **allow_deployments_to_no_targets** (Boolean, Optional)
- **exclude_unhealthy_targets** (Boolean, Optional)
- **skip_machine_behavior** (String, Optional)
- **target_roles** (List of String, Optional)


<a id="nestedblock--release_creation_strategy"></a>
### Nested Schema for `release_creation_strategy`

Optional:

- **channel_id** (String, Optional)
- **release_creation_package** (Block List, Max: 1) (see [below for nested schema](#nestedblock--release_creation_strategy--release_creation_package))
- **release_creation_package_step_id** (String, Optional)

<a id="nestedblock--release_creation_strategy--release_creation_package"></a>
### Nested Schema for `release_creation_strategy.release_creation_package`

Optional:

- **deployment_action** (String, Optional)
- **package_reference** (String, Optional)



<a id="nestedblock--templates"></a>
### Nested Schema for `templates`

Required:

- **name** (String, Required) The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`.

Optional:

- **default_value** (String, Optional) A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.
- **display_settings** (Map of String, Optional) The display settings for the parameter.
- **help_text** (String, Optional) The help presented alongside the parameter input.
- **id** (String, Optional) The unique ID for this resource.
- **label** (String, Optional) The label shown beside the parameter when presented in the deployment process. Example: `Server name`.


<a id="nestedblock--version_control_settings"></a>
### Nested Schema for `version_control_settings`

Optional:

- **password** (String, Optional) The password associated with this resource.
- **username** (String, Optional) The username associated with this resource.

Read-only:

- **default_branch** (String, Read-only)
- **url** (String, Read-only)


<a id="nestedblock--extension_settings"></a>
### Nested Schema for `extension_settings`

Read-only:

- **extension_id** (String, Read-only)
- **values** (List of String, Read-only)


<a id="nestedblock--versioning_strategy"></a>
### Nested Schema for `versioning_strategy`

Read-only:

- **donor_package** (List of String, Read-only)
- **donor_package_step_id** (String, Read-only)
- **template** (String, Read-only)


