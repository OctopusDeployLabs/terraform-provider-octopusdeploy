resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
  tenant_tags = ["tag1/a", "tag1/b"]
  depends_on = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
}

resource "octopusdeploy_tenant_project" "team_a_frontend" {
  tenant_id = octopusdeploy_tenant.tenant_team_a.id
  project_id = octopusdeploy_project.deploy_frontend_project.id
  environment_ids = [octopusdeploy_environment.test_environment.id, octopusdeploy_environment.development_environment.id, octopusdeploy_environment.production_environment.id]
}