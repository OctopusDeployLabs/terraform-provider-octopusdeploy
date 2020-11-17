provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_azure_service_principal" "azure_service_principal" {
  application_id       = var.application_id
  application_password = var.application_password
  key                  = var.key
  name                 = "azure account"
  subscription_id      = var.subscription_id
  tenant_id            = var.tenant_id
}
