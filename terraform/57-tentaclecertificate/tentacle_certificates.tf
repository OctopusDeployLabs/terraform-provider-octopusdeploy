resource "octopusdeploy_tentacle_certificate" "base" {}

resource "octopusdeploy_tentacle_certificate" "optional" {
  dependencies = {
    "base_id" = octopusdeploy_tentacle_certificate.base.id
  }
}

output "base_certificate_thumbprint" {
  value = octopusdeploy_tentacle_certificate.base.thumbprint
}