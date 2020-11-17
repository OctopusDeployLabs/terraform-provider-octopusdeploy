provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_project" "DevOpsProject" {
  lifecycle_id     = "Default Lifecycle"
  name             = "Terratest"
  project_group_id = "Projects-1"
}
