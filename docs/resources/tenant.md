---
page_title: "octopusdeploy_tenant Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages tenants in Octopus Deploy.
---

# octopusdeploy_tenant (Resource)

This resource manages tenants in Octopus Deploy.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of this resource.

### Optional

- `cloned_from_tenant_id` (String) The ID of the tenant from which this tenant was cloned.
- `description` (String) The description of this tenant.
- `is_disabled` (Boolean) The disabled status of this tenant.
- `space_id` (String) The space ID associated with this tenant.
- `tenant_tags` (Set of String) A list of tenant tags associated with this resource.

### Read-Only

- `id` (String) The unique ID for this resource.

~> **NOTE property `project_environment` deprecated:** The `project_environment` property has been replaced by the `octopusdeploy_tenant_project` resource to allow more advanced provisioning scenarioes. 