---
page_title: "octopusdeploy_channels Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing channels.
---

# Data Source `octopusdeploy_channels`

Provides information about existing channels.

## Example Usage

```terraform
data "octopusdeploy_channels" "example" {
  ids          = ["Channels-123", "Channels-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **channel** (Block List) A channel that matches the specified filter(s). (see [below for nested schema](#nestedblock--channel))

<a id="nestedblock--channel"></a>
### Nested Schema for `channel`

Read-only:

- **description** (String, Read-only) The description of this resource.
- **id** (String, Read-only) The unique ID for this resource.
- **is_default** (Boolean, Read-only) Indicates if this is the default channel for the associated project.
- **lifecycle_id** (String, Read-only) The lifecycle ID associated with this channel.
- **name** (String, Read-only) The name of this resource.
- **project_id** (String, Read-only) The project ID associated with this channel.
- **rules** (List of Object, Read-only) A list of rules associated with this channel. (see [below for nested schema](#nestedatt--channel--rules))
- **space_id** (String, Read-only) The space ID associated with this resource.
- **tenant_tags** (List of String, Read-only) A list of tenant tags associated with this resource.

<a id="nestedatt--channel--rules"></a>
### Nested Schema for `channel.rules`

- **actions** (List of String)
- **id** (String)
- **tag** (String)
- **version_range** (String)


