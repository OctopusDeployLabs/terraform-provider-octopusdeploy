---
page_title: "octopusdeploy_certificate Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages certificates in Octopus Deploy.
---

# Resource `octopusdeploy_certificate`

This resource manages certificates in Octopus Deploy.

## Example Usage

```terraform
resource "octopusdeploy_certificate" "example" {
  certificate_data = "a-base-64-encoded-string-representing-the-certificate-data"
  name             = "Development Certificate"
  password         = "some-random-value"
}
```

## Schema

### Required

- **certificate_data** (String, Required) The encoded data of the certificate.
- **name** (String, Required) The name of this resource.
- **password** (String, Required) The password associated with this resource.

### Optional

- **archived** (String, Optional)
- **certificate_data_format** (String, Optional) Specifies the archive file format used for storing cryptography objects in the certificate. Valid formats are `Der`, `Pem`, `Pkcs12`, or `Unknown`.
- **environments** (List of String, Optional) A list of environment IDs associated with this resource.
- **has_private_key** (Boolean, Optional) Indicates if the certificate has a private key.
- **id** (String, Optional) The unique ID for this resource.
- **is_expired** (Boolean, Optional) Indicates if the certificate has expired.
- **issuer_common_name** (String, Optional)
- **issuer_distinguished_name** (String, Optional)
- **issuer_organization** (String, Optional)
- **not_after** (String, Optional)
- **not_before** (String, Optional)
- **notes** (String, Optional)
- **replaced_by** (String, Optional)
- **self_signed** (Boolean, Optional)
- **serial_number** (String, Optional)
- **signature_algorithm_name** (String, Optional)
- **subject_alternative_names** (List of String, Optional)
- **subject_common_name** (String, Optional)
- **subject_distinguished_name** (String, Optional)
- **subject_organization** (String, Optional)
- **tenant_tags** (List of String, Optional) A list of tenant tags associated with this resource.
- **tenanted_deployment_participation** (String, Optional) The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.
- **tenants** (List of String, Optional) A list of tenant IDs associated with this resource.
- **thumbprint** (String, Optional)
- **version** (Number, Optional)


