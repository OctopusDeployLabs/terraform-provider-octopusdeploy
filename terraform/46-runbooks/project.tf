data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_project" "deploy_frontend_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test"
  project_group_id                     = octopusdeploy_project_group.project_group_test.id
  tenanted_deployment_participation    = "Untenanted"
  space_id                             = var.octopus_space_id
  included_library_variable_sets       = []
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_runbook" "runbook" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook"
  description        = "Test Runbook"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep    = 10
  }
  environment_scope           = "Specified"
  environments                = [octopusdeploy_environment.development_environment.id]
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "runbook" {
  runbook_id = octopusdeploy_runbook.runbook.id

  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Hello world (using PowerShell)"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = true
      worker_pool_id                     = ""
      properties                         = {
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.ScriptBody"   = "Write-Host 'Hello world, using PowerShell'\n\n#TODO: Experiment with steps of your own :)\n\nWrite-Host '[Learn more about the types of steps available in Octopus](https://oc.to/OnboardingAddStepsLearnMore)'"
        "Octopus.Action.Script.Syntax"       = "PowerShell"
      }
      environments          = [octopusdeploy_environment.development_environment.id]
      excluded_environments = []
      channels              = []
      tenant_tags           = []
      features              = []
    }

    properties   = {}
    target_roles = []
  }
}

resource "octopusdeploy_runbook" "runbook2" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook 2"
  description        = "Test Runbook 2"
  multi_tenancy_mode = "Tenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    should_keep_forever = true
  }
  environment_scope           = "All"
  environments                = []
  default_guided_failure_mode = "On"
  force_package_download      = false
}

resource "octopusdeploy_runbook" "runbook3" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook 3"
  description        = "Test Runbook 3"
  multi_tenancy_mode = "TenantedOrUntenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = true
    exclude_unhealthy_targets       = true
    skip_machine_behavior           = "None"
  }
  retention_policy {
    should_keep_forever = true
  }
  environment_scope           = "FromProjectLifecycles"
  environments                = []
  default_guided_failure_mode = "Off"
}

resource "octopusdeploy_runbook" "runbook4" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook 4"
}