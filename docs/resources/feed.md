---
page_title: "octopusdeploy_feed Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  
---

# Resource `octopusdeploy_feed`





## Schema

### Required

- **feed_uri** (String, Required)
- **name** (String, Required) The name of this resource.
- **package_acquisition_location_options** (List of String, Required)

### Optional

- **access_key** (String, Optional)
- **api_version** (String, Optional)
- **delete_unreleased_packages_after_days** (Number, Optional)
- **download_attempts** (Number, Optional)
- **download_retry_backoff_seconds** (Number, Optional)
- **feed_type** (String, Optional)
- **id** (String, Optional) The unique identifier for this resource.
- **is_enhanced_mode** (Boolean, Optional)
- **password** (String, Optional) The password associated with this resource.
- **registry_path** (String, Optional)
- **secret_key** (String, Optional)
- **space_id** (String, Optional) The space identifier associated with this resource.
- **username** (String, Optional) The username associated with this resource.

### Read-only

- **region** (String, Read-only)


