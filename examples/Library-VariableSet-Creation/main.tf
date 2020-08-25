provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_library_variable_set" "newaccount" {
  name            = var.variableSetName
  description     = var.description
}
