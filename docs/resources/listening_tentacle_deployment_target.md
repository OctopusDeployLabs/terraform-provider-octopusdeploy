---
page_title: "octopusdeploy_listening_tentacle_deployment_target Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages listening tentacle deployment targets in Octopus Deploy.
---

# Resource `octopusdeploy_listening_tentacle_deployment_target`

This resource manages listening tentacle deployment targets in Octopus Deploy.



## Schema

### Required

- **name** (String, Required) The name of this resource.
- **roles** (List of String, Required) A list of role IDs that are associated with this deployment target.
- **tentacle_url** (String, Required) The tenant URL of this deployment target.

### Optional

- **certificate_signature_algorithm** (String, Optional)
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **health_status** (String, Optional) Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.
- **id** (String, Optional) The unique ID for this resource.
- **is_disabled** (Boolean, Optional) Represents the disabled status of this deployment target.
- **is_in_process** (Boolean, Optional) Represents the in-process status of this deployment target.
- **machine_policy_id** (String, Optional) The machine policy ID that is associated with this deployment target.
- **operating_system** (String, Optional) The operating system that is associated with this deployment target.
- **proxy_id** (String, Optional) The proxy ID that is associated with this deployment target.
- **shell_name** (String, Optional) The shell name associated with this deployment target.
- **shell_version** (String, Optional) The shell version associated with this deployment target.
- **space_id** (String, Optional) The space ID associated with this resource.
- **status** (String, Optional) The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.
- **status_summary** (String, Optional) A summary elaborating on the status of this resource.
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.
- **tentacle_version_details** (Block List) (see [below for nested schema](#nestedblock--tentacle_version_details))
- **thumbprint** (String, Optional) The thumbprint of this deployment target.
- **uri** (String, Optional) The URI of this deployment target.

### Read-only

- **has_latest_calamari** (Boolean, Read-only)

<a id="nestedblock--tentacle_version_details"></a>
### Nested Schema for `tentacle_version_details`

Optional:

- **upgrade_locked** (Boolean, Optional)
- **upgrade_required** (Boolean, Optional)
- **upgrade_suggested** (Boolean, Optional)
- **version** (String, Optional)


