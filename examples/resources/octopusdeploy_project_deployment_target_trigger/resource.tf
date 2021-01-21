resource "octopusdeploy_project_deployment_target_trigger" "example" {
  name             = "[deployment_target_trigger_name]"
  project_id       = "Projects-123"
  event_categories = ["MachineUnhealthy"]
}