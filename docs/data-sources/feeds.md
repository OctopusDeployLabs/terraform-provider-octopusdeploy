---
page_title: "octopusdeploy_feeds Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing feeds.
---

# Data Source `octopusdeploy_feeds`

Provides information about existing feeds.



## Schema

### Optional

- **feed_type** (String, Optional) A filter to search by feed type. Valid feed types are `AwsElasticContainerRegistry`, `BuiltIn`, `Docker`, `GitHub`, `Helm`, `Maven`, `NuGet`, or `OctopusProject`.
- **id** (String, Optional) The ID of this resource.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.

### Read-only

- **feeds** (Block List) A list of feeds that match the filter(s). (see [below for nested schema](#nestedblock--feeds))

<a id="nestedblock--feeds"></a>
### Nested Schema for `feeds`

Read-only:

- **access_key** (String, Read-only)
- **api_version** (String, Read-only)
- **delete_unreleased_packages_after_days** (Number, Read-only)
- **download_attempts** (Number, Read-only)
- **download_retry_backoff_seconds** (Number, Read-only)
- **feed_type** (String, Read-only)
- **feed_uri** (String, Read-only)
- **id** (String, Read-only) The unique identifier for this resource.
- **is_enhanced_mode** (Boolean, Read-only)
- **name** (String, Read-only) The name of this resource.
- **package_acquisition_location_options** (List of String, Read-only)
- **password** (String, Read-only) The password associated with this resource.
- **region** (String, Read-only)
- **registry_path** (String, Read-only)
- **secret_key** (String, Read-only)
- **space_id** (String, Read-only) The space identifier associated with this resource.
- **username** (String, Read-only) The username associated with this resource.


