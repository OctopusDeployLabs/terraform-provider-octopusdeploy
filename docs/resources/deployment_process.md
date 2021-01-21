---
page_title: "octopusdeploy_deployment_process Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages deployment processes in Octopus Deploy.
---

# Resource `octopusdeploy_deployment_process`

This resource manages deployment processes in Octopus Deploy.



## Schema

### Required

- **project_id** (String, Required)

### Optional

- **id** (String, Optional) The unique ID for this resource.
- **last_snapshot_id** (String, Optional)
- **step** (Block List) (see [below for nested schema](#nestedblock--step))
- **version** (Number, Optional)

<a id="nestedblock--step"></a>
### Nested Schema for `step`

Required:

- **name** (String, Required) The name of this resource.

Optional:

- **action** (Block List) (see [below for nested schema](#nestedblock--step--action))
- **apply_terraform_action** (Block List) (see [below for nested schema](#nestedblock--step--apply_terraform_action))
- **condition** (String, Optional) When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'
- **condition_expression** (String, Optional) The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'
- **deploy_kubernetes_secret_action** (Block List) (see [below for nested schema](#nestedblock--step--deploy_kubernetes_secret_action))
- **deploy_package_action** (Block List) (see [below for nested schema](#nestedblock--step--deploy_package_action))
- **deploy_windows_service_action** (Block List) (see [below for nested schema](#nestedblock--step--deploy_windows_service_action))
- **id** (String, Optional) The unique ID for this resource.
- **manual_intervention_action** (Block List) (see [below for nested schema](#nestedblock--step--manual_intervention_action))
- **package_requirement** (String, Optional) Whether to run this step before or after package acquisition (if possible)
- **run_kubectl_script_action** (Block List) (see [below for nested schema](#nestedblock--step--run_kubectl_script_action))
- **run_script_action** (Block List) (see [below for nested schema](#nestedblock--step--run_script_action))
- **start_trigger** (String, Optional) Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')
- **target_roles** (List of String, Optional) The roles that this step run against, or runs on behalf of
- **window_size** (String, Optional) The maximum number of targets to deploy to simultaneously

<a id="nestedblock--step--action"></a>
### Nested Schema for `step.action`

Required:

- **action_type** (String, Required) The type of action
- **name** (String, Required) The name of this resource.

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **package** (Block Set) The primary package for the action (see [below for nested schema](#nestedblock--step--action--package))
- **primary_package** (Block Set, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--action--primary_package))
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **run_on_server** (Boolean, Optional) Whether this step runs on a worker or on the target
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **worker_pool_id** (String, Optional) Which worker pool to run on

<a id="nestedblock--step--action--package"></a>
### Nested Schema for `step.action.package`

Required:

- **name** (String, Required) The name of the package
- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **extract_during_deployment** (Boolean, Optional) Whether to extract the package during deployment
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--action--package--property))

<a id="nestedblock--step--action--package--property"></a>
### Nested Schema for `step.action.package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--action--primary_package"></a>
### Nested Schema for `step.action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--action--primary_package--property))

<a id="nestedblock--step--action--primary_package--property"></a>
### Nested Schema for `step.action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--action--property"></a>
### Nested Schema for `step.action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--apply_terraform_action"></a>
### Nested Schema for `step.apply_terraform_action`

Required:

- **name** (String, Required) The name of this resource.

Optional:

- **additional_init_params** (String, Optional) Additional parameters passed to the init command
- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **primary_package** (Block Set, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--apply_terraform_action--primary_package))
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--apply_terraform_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **run_on_server** (Boolean, Optional) Whether this step runs on a worker or on the target
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--step--apply_terraform_action--primary_package"></a>
### Nested Schema for `step.apply_terraform_action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--apply_terraform_action--primary_package--property))

<a id="nestedblock--step--apply_terraform_action--primary_package--property"></a>
### Nested Schema for `step.apply_terraform_action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--apply_terraform_action--property"></a>
### Nested Schema for `step.apply_terraform_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--deploy_kubernetes_secret_action"></a>
### Nested Schema for `step.deploy_kubernetes_secret_action`

Required:

- **name** (String, Required) The name of this resource.
- **secret_name** (String, Required) The name of the secret resource
- **secret_values** (Block List, Min: 1) (see [below for nested schema](#nestedblock--step--deploy_kubernetes_secret_action--secret_values))

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--deploy_kubernetes_secret_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **run_on_server** (Boolean, Optional) Whether this step runs on a worker or on the target
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--step--deploy_kubernetes_secret_action--secret_values"></a>
### Nested Schema for `step.deploy_kubernetes_secret_action.secret_values`

Required:

- **key** (String, Required)
- **value** (String, Required)


<a id="nestedblock--step--deploy_kubernetes_secret_action--property"></a>
### Nested Schema for `step.deploy_kubernetes_secret_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--deploy_package_action"></a>
### Nested Schema for `step.deploy_package_action`

Required:

- **name** (String, Required) The name of this resource.
- **primary_package** (Block Set, Min: 1, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--deploy_package_action--primary_package))

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--deploy_package_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **windows_service** (Block Set, Max: 1) Deploy a windows service feature (see [below for nested schema](#nestedblock--step--deploy_package_action--windows_service))

<a id="nestedblock--step--deploy_package_action--primary_package"></a>
### Nested Schema for `step.deploy_package_action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--deploy_package_action--primary_package--property))

<a id="nestedblock--step--deploy_package_action--primary_package--property"></a>
### Nested Schema for `step.deploy_package_action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--deploy_package_action--property"></a>
### Nested Schema for `step.deploy_package_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action


<a id="nestedblock--step--deploy_package_action--windows_service"></a>
### Nested Schema for `step.deploy_package_action.windows_service`

Required:

- **executable_path** (String, Required) The path to the executable relative to the package installation directory
- **service_name** (String, Required) The name of the service

Optional:

- **arguments** (String, Optional) The command line arguments that will be passed to the service when it starts
- **custom_account_name** (String, Optional) The Windows/domain account of the custom user that the service will run under
- **custom_account_password** (String, Optional) The password for the custom account
- **dependencies** (String, Optional) Any dependencies that the service has. Separate the names using forward slashes (/).
- **description** (String, Optional) User-friendly description of the service (optional)
- **display_name** (String, Optional) The display name of the service (optional)
- **service_account** (String, Optional) Which built-in account will the service run under. Can be LocalSystem, NT Authority\NetworkService, NT Authority\LocalService, _CUSTOM or an expression
- **start_mode** (String, Optional) When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression



<a id="nestedblock--step--deploy_windows_service_action"></a>
### Nested Schema for `step.deploy_windows_service_action`

Required:

- **executable_path** (String, Required) The path to the executable relative to the package installation directory
- **name** (String, Required) The name of this resource.
- **primary_package** (Block Set, Min: 1, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--deploy_windows_service_action--primary_package))
- **service_name** (String, Required) The name of the service

Optional:

- **arguments** (String, Optional) The command line arguments that will be passed to the service when it starts
- **channels** (List of String, Optional) The channels that this step applies to
- **custom_account_name** (String, Optional) The Windows/domain account of the custom user that the service will run under
- **custom_account_password** (String, Optional) The password for the custom account
- **dependencies** (String, Optional) Any dependencies that the service has. Separate the names using forward slashes (/).
- **description** (String, Optional) User-friendly description of the service (optional)
- **disabled** (Boolean, Optional) Whether this step is disabled
- **display_name** (String, Optional) The display name of the service (optional)
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--deploy_windows_service_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **service_account** (String, Optional) Which built-in account will the service run under. Can be LocalSystem, NT Authority\NetworkService, NT Authority\LocalService, _CUSTOM or an expression
- **start_mode** (String, Optional) When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--step--deploy_windows_service_action--primary_package"></a>
### Nested Schema for `step.deploy_windows_service_action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--deploy_windows_service_action--primary_package--property))

<a id="nestedblock--step--deploy_windows_service_action--primary_package--property"></a>
### Nested Schema for `step.deploy_windows_service_action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--deploy_windows_service_action--property"></a>
### Nested Schema for `step.deploy_windows_service_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--manual_intervention_action"></a>
### Nested Schema for `step.manual_intervention_action`

Required:

- **instructions** (String, Required) The instructions for the user to follow
- **name** (String, Required) The name of this resource.

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--manual_intervention_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **responsible_teams** (String, Optional) The teams responsible to resolve this step. If no teams are specified, all users who have permission to deploy the project can resolve it.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--step--manual_intervention_action--property"></a>
### Nested Schema for `step.manual_intervention_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_kubectl_script_action"></a>
### Nested Schema for `step.run_kubectl_script_action`

Required:

- **name** (String, Required) The name of this resource.

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **package** (Block Set) The primary package for the action (see [below for nested schema](#nestedblock--step--run_kubectl_script_action--package))
- **primary_package** (Block Set, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--run_kubectl_script_action--primary_package))
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_kubectl_script_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **run_on_server** (Boolean, Optional) Whether this step runs on a worker or on the target
- **script_file_name** (String, Optional) The script file name in the package
- **script_parameters** (String, Optional) Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.

<a id="nestedblock--step--run_kubectl_script_action--package"></a>
### Nested Schema for `step.run_kubectl_script_action.package`

Required:

- **name** (String, Required) The name of the package
- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **extract_during_deployment** (Boolean, Optional) Whether to extract the package during deployment
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_kubectl_script_action--package--property))

<a id="nestedblock--step--run_kubectl_script_action--package--property"></a>
### Nested Schema for `step.run_kubectl_script_action.package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_kubectl_script_action--primary_package"></a>
### Nested Schema for `step.run_kubectl_script_action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_kubectl_script_action--primary_package--property))

<a id="nestedblock--step--run_kubectl_script_action--primary_package--property"></a>
### Nested Schema for `step.run_kubectl_script_action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_kubectl_script_action--property"></a>
### Nested Schema for `step.run_kubectl_script_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_script_action"></a>
### Nested Schema for `step.run_script_action`

Required:

- **name** (String, Required) The name of this resource.

Optional:

- **channels** (List of String, Optional) The channels that this step applies to
- **disabled** (Boolean, Optional) Whether this step is disabled
- **environments** (List of String, Optional) The environments that this step will run in
- **excluded_environments** (List of String, Optional) The environments that this step will be skipped in
- **package** (Block Set) The primary package for the action (see [below for nested schema](#nestedblock--step--run_script_action--package))
- **primary_package** (Block Set, Max: 1) The primary package for the action (see [below for nested schema](#nestedblock--step--run_script_action--primary_package))
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_script_action--property))
- **required** (Boolean, Optional) Whether this step is required and cannot be skipped
- **run_on_server** (Boolean, Optional) Whether this step runs on a worker or on the target
- **script_file_name** (String, Optional) The script file name in the package
- **script_parameters** (String, Optional) Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **variable_substitution_in_files** (String, Optional) A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.

<a id="nestedblock--step--run_script_action--package"></a>
### Nested Schema for `step.run_script_action.package`

Required:

- **name** (String, Required) The name of the package
- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **extract_during_deployment** (Boolean, Optional) Whether to extract the package during deployment
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_script_action--package--property))

<a id="nestedblock--step--run_script_action--package--property"></a>
### Nested Schema for `step.run_script_action.package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_script_action--primary_package"></a>
### Nested Schema for `step.run_script_action.primary_package`

Required:

- **package_id** (String, Required) The ID of the package

Optional:

- **acquisition_location** (String, Optional) Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression
- **feed_id** (String, Optional) The feed to retrieve the package from
- **property** (Block Set) (see [below for nested schema](#nestedblock--step--run_script_action--primary_package--property))

<a id="nestedblock--step--run_script_action--primary_package--property"></a>
### Nested Schema for `step.run_script_action.primary_package.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action



<a id="nestedblock--step--run_script_action--property"></a>
### Nested Schema for `step.run_script_action.property`

Required:

- **key** (String, Required) The name of the action
- **value** (String, Required) The type of action


