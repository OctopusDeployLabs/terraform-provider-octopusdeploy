resource "octopusdeploy_space" "example" {
  description                 = "A space for the development team."
  name                        = "Development Team Space"
  is_default                  = false
  is_task_queue_stopped       = false
  space_managers_team_members = ["Users-123", "Users-321"]
  space_managers_teams        = ["teams-everyone"]
}