provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_project" "DevOpsProject" {
  name             = "testProject"
  lifecycle_id     = "Default Lifecycle"
  project_group_id = "Dev"
}
