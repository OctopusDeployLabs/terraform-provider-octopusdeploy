provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space_id   = var.space
}

resource "octopusdeploy_certificate" "DevCert" {
  name = "DevCert"

  certificate_data = var.certEOM
}