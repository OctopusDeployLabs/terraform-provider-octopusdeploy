# Deployment freeze
resource "octopusdeploy_deployment_freeze" "freeze" {
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
}


resource "octopusdeploy_deployment_freeze_project" "project_freeze" {
  deploymentfreeze_id= octopusdeploy_deployment_freeze.freeze.id
  project_id = "Projects-123"
  environment_ids = [ "Environments-123", "Environments-456" ]
}

# Freeze with ids sourced from resources and data sources. 
# Projects can be sourced from different spaces, a single scope can only reference projects and environments from the same space.
resource "octopusdeploy_deployment_freeze" "freeze" {
  name = "End of financial year shutdown"
  start = "2025-06-30T00:00:00+10:00"
  end = "2025-07-02T00:00:00+10:00"
}

resource "octopusdeploy_deployment_freeze_project" "project_freeze" {
  deploymentfreeze_id = octopusdeploy_deployment_freeze.freeze.id
  project_id          = resource.octopusdeploy_project.project1.id
  environment_ids = [resource.octopusdeploy_environment.production.id]
}

resource "octopusdeploy_deployment_freeze_project" "project_freeze" {
  deploymentfreeze_id = octopusdeploy_deployment_freeze.freeze.id
  project_id          = data.octopusdeploy_projects.second_project.projects[0].id
  environment_ids = [ data.octopusdeploy_environments.default_environment.environments[0].id ]
}