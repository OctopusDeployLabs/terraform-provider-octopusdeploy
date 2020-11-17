provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_certificate" "developer_certificate" {
  certificate_data = var.certEOM
  name             = "developer certificate"
}
