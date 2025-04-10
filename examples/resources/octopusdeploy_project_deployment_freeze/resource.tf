# basic freeze with no environment scopes
resource "octopusdeploy_project_deployment_freeze" "freeze" {
  owner_id = "Projects-123"
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
}

# Freeze with different timezones and single environment scope
resource "octopusdeploy_deployment_freeze" "freeze" {
  owner_id = "Projects-123"
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
  environment_ids = ["Environments-123"]
}

# Freeze recurring freeze yearly on Xmas
resource "octopusdeploy_deployment_freeze" "freeze" {
  owner_id = "Projects-123"
  name = "Xmas"
  start = "2024-12-25T00:00:00+10:00"
  end = "2024-12-27T00:00:00+08:00"
  recurring_schedule = {
    type    = "Annually"
    unit    = 1
    end_type = "Never"
  }
}

resource "octopusdeploy_deployment_freeze_project" "project_freeze" {
  deploymentfreeze_id= octopusdeploy_deployment_freeze.freeze.id
  project_id = "Projects-123"
  environment_ids = [ "Environments-123", "Environments-456" ]
}

# Freeze with ids sourced from resources.
resource "octopusdeploy_deployment_freeze" "freeze" {
  owner_id = resource.octopusdeploy_project.project1.id
  name = "End of financial year shutdown"
  start = "2025-06-30T00:00:00+10:00"
  end = "2025-07-02T00:00:00+10:00"
  environment_ids = [resource.octopusdeploy_environment.production.id]
}

# Freeze with tenant environment scope and ids sourced from datasources.
resource "octopusdeploy_deployment_freeze" "freeze" {
  owner_id = data.octopusdeploy_project.project1.id
  name = "End of financial year shutdown"
  start = "2025-06-30T00:00:00+10:00"
  end = "2025-07-02T00:00:00+10:00"
  environment_ids = [data.octopusdeploy_environments.default_environment.environments[0].id]
}

resource "octopusdeploy_deployment_freeze_tenant" "tenant_freeze" {
  deploymentfreeze_id = octopusdeploy_deployment_freeze.freeze.id
  tenant_id           = data.octopusdeploy_tenants.default_tenant.tenants[0].id
  project_id          = data.octopusdeploy_project.project1.id
  environment_id      = data.octopusdeploy_environments.default_environment.environments[0].id
}