resource "octopusdeploy_variable" "unscoped_project_variable" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "UnscopedVariable"
  value    = "UnscopedVariable"
  depends_on = [octopusdeploy_project.test_project]
}

resource "octopusdeploy_variable" "scoped_project_variable_action" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ActionScopedVariable"
  value    = "unscoped variable"
  scope {
    actions = [octopusdeploy_deployment_process.test_deployment_process.step[0].run_script_action[0].id]
  }
  depends_on = [octopusdeploy_project.test_project, octopusdeploy_deployment_process.test_deployment_process]
}

resource "octopusdeploy_variable" "scoped_project_variable_channel" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ChannelScopedVariable"
  value    = "ChannelScopedVariable"
  scope {
    channels = [octopusdeploy_channel.test_channel.id]
  }
  depends_on = [octopusdeploy_project.test_project, octopusdeploy_channel.test_channel]
}

resource "octopusdeploy_variable" "scoped_project_variable_environment" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "EnvironmentScopedVariable"
  value    = "EnvironmentScopedVariable"
  scope {
    environments = [octopusdeploy_environment.development_environment.id]
  }
  depends_on = [octopusdeploy_project.test_project, octopusdeploy_environment.development_environment]
}

resource "octopusdeploy_variable" "scoped_project_variable_machine" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "MachineScopedVariable"
  value    = "MachineScopedVariable"
  scope {
    machines = [octopusdeploy_cloud_region_deployment_target.test_target.id]
  }
  depends_on = [octopusdeploy_project.test_project, octopusdeploy_cloud_region_deployment_target.test_target]
}

resource "octopusdeploy_variable" "scoped_project_variable_process" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ProcessScopedVariable"
  value    = "ProcessScopedVariable"
  scope {
    processes = [octopusdeploy_deployment_process.test_deployment_process.id]
  }
  depends_on = [octopusdeploy_project.test_project, octopusdeploy_deployment_process.test_deployment_process]
}

resource "octopusdeploy_variable" "scoped_project_variable_role" {
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "RoleScopedVariable"
  value    = "RoleScopedVariable"
  scope {
    roles = ["role"]
  }
  depends_on = [octopusdeploy_project.test_project]
}