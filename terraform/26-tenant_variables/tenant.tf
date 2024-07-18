resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
}

resource "octopusdeploy_tenant_project" "team_a_frontend" {
  tenant_id = octopusdeploy_tenant.tenant_team_a.id
  project_id = octopusdeploy_project.deploy_frontend_project.id
  environment_ids = [octopusdeploy_environment.test_environment.id, octopusdeploy_environment.development_environment.id, octopusdeploy_environment.production_environment.id]
}