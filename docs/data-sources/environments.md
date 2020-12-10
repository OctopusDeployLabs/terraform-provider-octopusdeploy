---
page_title: "octopusdeploy_environments Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing environments.
---

# Data Source `octopusdeploy_environments`

Provides information about existing environments.

## Example Usage

```terraform
data "octopusdeploy_environments" "example" {
  ids          = ["Environments-123", "Environments-321"]
  name         = "Production"
  partial_name = "Produc"
  skip         = 5
  take         = 100
}
```

## Schema

### Optional

- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **name** (String, Optional) A filter to search by name.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **environment** (Block List) A list of environments that match the filter(s). (see [below for nested schema](#nestedblock--environment))
- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedblock--environment"></a>
### Nested Schema for `environment`

Read-only:

- **allow_dynamic_infrastructure** (Boolean, Read-only)
- **description** (String, Read-only) The description of this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **sort_order** (Number, Read-only)
- **use_guided_failure** (Boolean, Read-only)


