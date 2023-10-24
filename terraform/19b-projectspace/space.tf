resource "octopusdeploy_space" "octopus_project_space_test" {
  name                  = "Project Space Test"
  is_default            = false
  is_task_queue_stopped = false
  description           = "My test space"
  space_managers_teams  = ["teams-administrators"]
}
