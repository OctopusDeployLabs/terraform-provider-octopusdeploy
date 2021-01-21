provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_library_variable_set" "newaccount" {
  description = var.description
  name        = var.variableSetName
}
