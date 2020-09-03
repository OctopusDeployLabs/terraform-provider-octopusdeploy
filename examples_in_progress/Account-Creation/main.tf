provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_account" "newaccount" {
  name            = var.azureAccountName
  account_type    = "Azure"
  subscription_id = var.subID
}
