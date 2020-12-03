---
page_title: "octopusdeploy_project_deployment_target_trigger Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_project_deployment_target_trigger`





## Schema

### Required

- **name** (String, Required) The name of this resource.
- **project_id** (String, Required) The project_id of the Project to attach the trigger to.

### Optional

- **environment_ids** (List of String, Optional) Apply environment id filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.
- **event_categories** (List of String, Optional) Apply event category filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.
- **event_groups** (List of String, Optional) Apply event group filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.
- **id** (String, Optional) The ID of this resource.
- **roles** (List of String, Optional) Apply event role filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.
- **should_redeploy** (Boolean, Optional) Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.


