---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_machine_proxies Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing Octopus Deploy machine proxies.
---

# octopusdeploy_machine_proxies (Data Source)

Provides information about existing Octopus Deploy machine proxies.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ids` (List of String) A filter to search by a list of IDs.
- `partial_name` (String) A filter to search by a partial name.
- `skip` (Number) A filter to specify the number of items to skip in the response.
- `space_id` (String) A Space ID to filter by. Will revert what is specified on the provider if not set
- `take` (Number) A filter to specify the number of items to take (or return) in the response.

### Read-Only

- `id` (String) An auto-generated identifier that includes the timestamp when this data source was last modified.
- `machine_proxies` (Attributes List) A list of machine proxies that match the filter(s). (see [below for nested schema](#nestedatt--machine_proxies))

<a id="nestedatt--machine_proxies"></a>
### Nested Schema for `machine_proxies`

Read-Only:

- `host` (String) DNS hostname of the proxy server
- `id` (String)
- `name` (String)
- `port` (Number) The port number for the proxy server.
- `space_id` (String) The space ID associated with this machine proxy.
- `username` (String) Username for the proxy server


