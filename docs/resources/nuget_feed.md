---
page_title: "octopusdeploy_nuget_feed Resource - terraform-provider-octopusdeploy"
subcategory: "Feeds"
description: |-
  This resource manages a Nuget feed in Octopus Deploy.
---

# octopusdeploy_nuget_feed (Resource)

This resource manages a Nuget feed in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_nuget_feed" "example" {
  download_attempts              = 1
  download_retry_backoff_seconds = 30
  feed_uri                       = "https://api.nuget.org/v3/index.json"
  is_enhanced_mode               = true
  password                       = "test-password"
  name                           = "Test NuGet Feed (OK to Delete)"
  username                       = "test-username"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `feed_uri` (String)
- `name` (String) The name of this resource.

### Optional

- `download_attempts` (Number) The number of times a deployment should attempt to download a package from this feed before failing.
- `download_retry_backoff_seconds` (Number) The number of seconds to apply as a linear back off between download attempts.
- `is_enhanced_mode` (Boolean) This will improve performance of the NuGet feed but may not be supported by some older feeds. Disable if the operation, Create Release does not return the latest version for a package.
- `package_acquisition_location_options` (List of String)
- `password` (String, Sensitive) The password associated with this resource.
- `space_id` (String) The space ID associated with this nuget feed.
- `username` (String, Sensitive) The username associated with this resource.

### Read-Only

- `id` (String) The unique ID for this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import [options] octopusdeploy_nuget_feed.<name> <feed-id>
```
