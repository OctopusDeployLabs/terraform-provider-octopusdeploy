---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_static_worker_pool Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages static worker pools in Octopus Deploy.
---

# octopusdeploy_static_worker_pool (Resource)

This resource manages static worker pools in Octopus Deploy.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of this resource.

### Optional

- **description** (String) The description of this resource.
- **id** (String) The unique ID for this resource.
- **is_default** (Boolean)
- **sort_order** (Number) The order number to sort a dynamic worker pool.
- **space_id** (String) The space ID associated with this resource.

### Read-Only

- **can_add_workers** (Boolean)

