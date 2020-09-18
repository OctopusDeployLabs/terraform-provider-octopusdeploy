provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space_id   = var.space
}

resource "octopusdeploy_azure_service_principal" "Azure" {
  name = "terratesttest"
  client_id = var.client_id
  tenant_id = var.tenant_id
  subscription_number = var.subscription_number
  key = var.key
}
