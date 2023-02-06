resource "octopusdeploy_project_deployment_target_trigger" "projecttrigger_test1" {
  name             = "Test 1"
  project_id       = "${octopusdeploy_project.project_1.id}"
  event_categories = []
  environment_ids  = []
  event_groups     = ["MachineAvailableForDeployment"]
  roles            = []
  should_redeploy  = false
}

resource "octopusdeploy_project_deployment_target_trigger" "projecttrigger_test2" {
  name             = "Test 2"
  project_id       = "${octopusdeploy_project.project_2.id}"
  event_categories = []
  environment_ids  = []
  event_groups     = ["MachineAvailableForDeployment"]
  roles            = []
  should_redeploy  = false
}