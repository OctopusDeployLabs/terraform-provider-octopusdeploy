---
page_title: "octopusdeploy_certificates Data Source - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  Provides information about existing certificates.
---

# Data Source `octopusdeploy_certificates`

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

## Schema

### Optional

- **archived** (String, Optional) A filter to search for resources that have been archived.
- **first_result** (String, Optional) A filter to define the first result.
- **ids** (List of String, Optional) A filter to search by a list of IDs.
- **order_by** (String, Optional) A filter used to order the search results.
- **partial_name** (String, Optional) A filter to search by the partial match of a name.
- **search** (String, Optional) A filter of terms used the search operation.
- **skip** (Number, Optional) A filter to specify the number of items to skip in the response.
- **take** (Number, Optional) A filter to specify the number of items to take (or return) in the response.
- **tenant** (String, Optional) A filter to search by a tenant ID.

### Read-only

- **certificates** (Block List) A list of certificates that match the filter(s). (see [below for nested schema](#nestedblock--certificates))
- **id** (String, Read-only) A auto-generated identifier that includes the timestamp when this data source was last modified.

<a id="nestedblock--certificates"></a>
### Nested Schema for `certificates`

Read-only:

- **archived** (String, Read-only)
- **certificate_data** (String, Read-only) The encoded data of the certificate.
- **certificate_data_format** (String, Read-only) Specifies the archive file format used for storing cryptography objects in the certificate. Valid formats are `Der`, `Pem`, `Pkcs12`, or `Unknown`.
- **environments** (List of String, Read-only) A list of environment IDs associated with this resource.
- **has_private_key** (Boolean, Read-only) Indicates if the certificate has a private key.
- **id** (String, Read-only) The unique ID for this resource.
- **is_expired** (Boolean, Read-only) Indicates if the certificate has expired.
- **issuer_common_name** (String, Read-only)
- **issuer_distinguished_name** (String, Read-only)
- **issuer_organization** (String, Read-only)
- **name** (String, Read-only) The name of this resource.
- **not_after** (String, Read-only)
- **not_before** (String, Read-only)
- **notes** (String, Read-only)
- **password** (String, Read-only) The password associated with this resource.
- **replaced_by** (String, Read-only)
- **self_signed** (Boolean, Read-only)
- **serial_number** (String, Read-only)
- **signature_algorithm_name** (String, Read-only)
- **subject_alternative_names** (List of String, Read-only)
- **subject_common_name** (String, Read-only)
- **subject_distinguished_name** (String, Read-only)
- **subject_organization** (String, Read-only)
- **tenant_tags** (List of String, Read-only) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Read-only) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Read-only) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Read-only)
- **version** (Number, Read-only)


