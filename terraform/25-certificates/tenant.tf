resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
  tenant_tags = ["tag1/a", "tag1/b"]
  depends_on = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
}



resource "octopusdeploy_tenant_project_environment" "team_a_test_environment" {
  tenant_id = octopusdeploy_tenant.tenant_team_a.id
  project_id = octopusdeploy_project.deploy_frontend_project.id
  environment_id = octopusdeploy_environment.test_environment.id
}

resource "octopusdeploy_tenant_project_environment" "team_a_development_environment" {
  tenant_id = octopusdeploy_tenant.tenant_team_a.id
  project_id = octopusdeploy_project.deploy_frontend_project.id
  environment_id = octopusdeploy_environment.development_environment.id
}

resource "octopusdeploy_tenant_project_environment" "team_a_production_environment" {
  tenant_id = octopusdeploy_tenant.tenant_team_a.id
  project_id = octopusdeploy_project.deploy_frontend_project.id
  environment_id = octopusdeploy_environment.production_environment.id
}