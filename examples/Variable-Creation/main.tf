provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_variable" "newVariable" {
  name       = var.varName
  project_id = var.projectID
  type       = "String"
}
