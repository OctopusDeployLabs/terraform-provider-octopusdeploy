# Deployment process
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
}

# Runbook process
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
  runbook_id  = "Runbooks-42"
}
