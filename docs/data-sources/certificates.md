---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_certificates Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing certificates.
---

# octopusdeploy_certificates (Data Source)

Provides information about existing certificates.

## Example Usage

```terraform
data "octopusdeploy_certificates" "example" {
  archived     = false
  ids          = ["Certificates-123", "Certificates-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `archived` (String) A filter to search for resources that have been archived.
- `first_result` (String) A filter to define the first result.
- `ids` (List of String) A filter to search by a list of IDs.
- `order_by` (String) A filter used to order the search results.
- `partial_name` (String) A filter to search by the partial match of a name.
- `search` (String) A filter of terms used the search operation.
- `skip` (Number) A filter to specify the number of items to skip in the response.
- `space_id` (String) The space ID associated with this resource.
- `take` (Number) A filter to specify the number of items to take (or return) in the response.
- `tenant` (String) A filter to search by a tenant ID.

### Read-Only

- `certificates` (List of Object) A list of certificates that match the filter(s). (see [below for nested schema](#nestedatt--certificates))
- `id` (String) An auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedatt--certificates"></a>
### Nested Schema for `certificates`

Read-Only:

- `archived` (String)
- `certificate_data` (String)
- `certificate_data_format` (String)
- `environments` (List of String)
- `has_private_key` (Boolean)
- `id` (String)
- `is_expired` (Boolean)
- `issuer_common_name` (String)
- `issuer_distinguished_name` (String)
- `issuer_organization` (String)
- `name` (String)
- `not_after` (String)
- `not_before` (String)
- `notes` (String)
- `password` (String)
- `replaced_by` (String)
- `self_signed` (Boolean)
- `serial_number` (String)
- `signature_algorithm_name` (String)
- `space_id` (String)
- `subject_alternative_names` (List of String)
- `subject_common_name` (String)
- `subject_distinguished_name` (String)
- `subject_organization` (String)
- `tenant_tags` (List of String)
- `tenanted_deployment_participation` (String)
- `tenants` (List of String)
- `thumbprint` (String)
- `version` (Number)


