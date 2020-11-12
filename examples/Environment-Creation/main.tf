provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_environment" "newEnvironment" {
  name = var.environmentName
  # Optional Inputs:
  # description                  = var.description
  # use_guided_failure           = var.use_guided_failure
  # allow_dynamic_infrastructure = var.allow_dynamic_infrastructure
  # sort_order                   = var.sort_order
}
