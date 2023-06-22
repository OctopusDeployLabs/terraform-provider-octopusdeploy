resource "octopusdeploy_tenant_project_variable" "tenantprojectvariable6_team_a" {
  environment_id = "${octopusdeploy_environment.development_environment.id}"
  project_id     = "${octopusdeploy_project.deploy_frontend_project.id}"
  template_id    = "${octopusdeploy_project.deploy_frontend_project.template[0].id}"
  tenant_id      = "${octopusdeploy_tenant.tenant_team_a.id}"
  value          = "my value"
}

resource "octopusdeploy_tenant_common_variable" "tenantcommonvariable1_team_a" {
  library_variable_set_id = "${octopusdeploy_library_variable_set.library_variable_set_octopus_variables.id}"
  template_id             = "${octopusdeploy_library_variable_set.library_variable_set_octopus_variables.template[0].id}"
  tenant_id               = "${octopusdeploy_tenant.tenant_team_a.id}"
  value                   = "my value"
}
