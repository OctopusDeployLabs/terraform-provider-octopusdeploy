---
page_title: "octopusdeploy_tag_sets Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing tag sets.
---

# Data Source `octopusdeploy_tag_sets`

Provides information about existing tag sets.



## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **tag_set** (Block List) A list of tag sets that match the filter(s). (see [below for nested schema](#nestedblock--tag_set))

<a id="nestedblock--tag_set"></a>
### Nested Schema for `tag_set`

Read-only:

- **id** (String, Read-only) The unique ID for this resource.
- **name** (String, Read-only) The name of this resource.
- **tags** (List of Object, Read-only) (see [below for nested schema](#nestedatt--tag_set--tags))

<a id="nestedatt--tag_set--tags"></a>
### Nested Schema for `tag_set.tags`

- **canonical_tag_name** (String)
- **color** (String)
- **description** (String)
- **id** (String)
- **name** (String)
- **sort_order** (Number)


