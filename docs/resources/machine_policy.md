---
page_title: "octopusdeploy_machine_policy Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages machine policies in Octopus Deploy.
---

# Resource `octopusdeploy_machine_policy`

This resource manages machine policies in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.

### Optional

- **connection_connect_timeout** (Number, Optional)
- **connection_retry_count_limit** (Number, Optional)
- **connection_retry_sleep_interval** (Number, Optional)
- **connection_retry_time_limit** (Number, Optional)
- **description** (String, Optional) The description of this resource.
- **id** (String, Optional) The unique ID for this resource.
- **machine_cleanup_policy** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--machine_cleanup_policy))
- **machine_connectivity_policy** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--machine_connectivity_policy))
- **machine_health_check_policy** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--machine_health_check_policy))
- **machine_update_policy** (Block Set, Max: 1) (see [below for nested schema](#nestedblock--machine_update_policy))
- **polling_request_maximum_message_processing_timeout** (Number, Optional)
- **polling_request_queue_timeout** (Number, Optional)
- **space_id** (String, Optional) The space ID associated with this resource.

### Read-only

- **is_default** (Boolean, Read-only)

<a id="nestedblock--machine_cleanup_policy"></a>
### Nested Schema for `machine_cleanup_policy`

Optional:

- **delete_machines_behavior** (String, Optional)
- **delete_machines_elapsed_timespan** (Number, Optional)


<a id="nestedblock--machine_connectivity_policy"></a>
### Nested Schema for `machine_connectivity_policy`

Optional:

- **machine_connectivity_behavior** (String, Optional)


<a id="nestedblock--machine_health_check_policy"></a>
### Nested Schema for `machine_health_check_policy`

Optional:

- **bash_health_check_policy** (Block List, Max: 1) (see [below for nested schema](#nestedblock--machine_health_check_policy--bash_health_check_policy))
- **health_check_cron** (String, Optional)
- **health_check_cron_timezone** (String, Optional)
- **health_check_interval** (Number, Optional)
- **health_check_type** (String, Optional)
- **powershell_health_check_policy** (Block List, Max: 1) (see [below for nested schema](#nestedblock--machine_health_check_policy--powershell_health_check_policy))

<a id="nestedblock--machine_health_check_policy--bash_health_check_policy"></a>
### Nested Schema for `machine_health_check_policy.bash_health_check_policy`

Optional:

- **run_type** (String, Optional)
- **script_body** (String, Optional)


<a id="nestedblock--machine_health_check_policy--powershell_health_check_policy"></a>
### Nested Schema for `machine_health_check_policy.powershell_health_check_policy`

Optional:

- **run_type** (String, Optional)
- **script_body** (String, Optional)



<a id="nestedblock--machine_update_policy"></a>
### Nested Schema for `machine_update_policy`

Optional:

- **calamari_update_behavior** (String, Optional)
- **tentacle_update_account_id** (String, Optional)
- **tentacle_update_behavior** (String, Optional)


