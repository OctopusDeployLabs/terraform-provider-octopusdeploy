provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space_id   = var.space
}

resource "octopusdeploy_project" "DevOpsProject" {
  name             = "Terratest"
  lifecycle_id     = "Default Lifecycle"
  project_group_id = "Projects-1"
}
