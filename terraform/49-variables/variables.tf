

resource "octopusdeploy_variable" "scoped_project_variable_action" {
  depends_on = [
    octopusdeploy_project.test_project,
  ]
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ActionScopedVariable"
  value    = "unscoped variable"
  scope {
    actions = [octopusdeploy_deployment_process.test_deployment_process.step[0].run_script_action[0].id]
  }
}

resource "octopusdeploy_variable" "unscoped_project_variable" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.scoped_project_variable_action,
  ]
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "UnscopedVariable"
  value    = "UnscopedVariable"
}

resource "octopusdeploy_variable" "scoped_project_variable_channel" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.unscoped_project_variable,
    octopusdeploy_variable.scoped_project_variable_action,
  ]
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ChannelScopedVariable"
  value    = "ChannelScopedVariable"
  scope {
    channels = [octopusdeploy_channel.test_channel.id]
  }
}

resource "octopusdeploy_variable" "scoped_project_variable_environment" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.unscoped_project_variable,
    octopusdeploy_variable.scoped_project_variable_action,
    octopusdeploy_variable.scoped_project_variable_channel,
  ]
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "EnvironmentScopedVariable"
  value    = "EnvironmentScopedVariable"
  scope {
    environments = [octopusdeploy_environment.development_environment.id]
  }
}

resource "octopusdeploy_variable" "scoped_project_variable_machine" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.unscoped_project_variable,
    octopusdeploy_variable.scoped_project_variable_action,
    octopusdeploy_variable.scoped_project_variable_channel,
    octopusdeploy_variable.scoped_project_variable_environment
  ]

  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "MachineScopedVariable"
  value    = "MachineScopedVariable"
  scope {
    machines = [octopusdeploy_cloud_region_deployment_target.test_target.id]
  }
}

resource "octopusdeploy_variable" "scoped_project_variable_process" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.unscoped_project_variable,
    octopusdeploy_variable.scoped_project_variable_action,
    octopusdeploy_variable.scoped_project_variable_channel,
    octopusdeploy_variable.scoped_project_variable_environment,
    octopusdeploy_variable.scoped_project_variable_machine,
  ]
  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "ProcessScopedVariable"
  value    = "ProcessScopedVariable"
  scope {
    processes = [octopusdeploy_deployment_process.test_deployment_process.id]
  }
}

resource "octopusdeploy_variable" "scoped_project_variable_role" {
  depends_on = [
    octopusdeploy_project.test_project,
    octopusdeploy_variable.unscoped_project_variable,
    octopusdeploy_variable.scoped_project_variable_action,
    octopusdeploy_variable.scoped_project_variable_channel,
    octopusdeploy_variable.scoped_project_variable_environment,
    octopusdeploy_variable.scoped_project_variable_machine,
    octopusdeploy_variable.scoped_project_variable_process,
  ]

  owner_id = octopusdeploy_project.test_project.id
  type     = "String"
  name     = "RoleScopedVariable"
  value    = "RoleScopedVariable"
  scope {
    roles = ["role"]
  }
}