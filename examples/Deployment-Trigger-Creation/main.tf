provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_project_deployment_target_trigger" "deploymentTrigger" {
    name = var.name
    project_id = var.projectID
    event_categories = ["MachineUnhealthy"]
}