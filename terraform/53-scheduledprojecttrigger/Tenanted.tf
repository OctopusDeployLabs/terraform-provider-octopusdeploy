resource "octopusdeploy_project" "tenanted" {
  space_id                             = var.octopus_space_id
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Tenanted"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Tenanted"
  project_group_id                     = octopusdeploy_project_group.project_group_test.id
  tenanted_deployment_participation    = "Tenanted"
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


resource "octopusdeploy_deployment_process" "tenanted_deployment_process" {
  project_id = octopusdeploy_project.tenanted.id
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
      run_on_server                      = true
    }
  }
}

resource "octopusdeploy_runbook" "tenanted_runbook" {
  project_id         = octopusdeploy_project.tenanted.id
  name               = "Runbook"
  description        = "Test Runbook"
  multi_tenancy_mode = "Tenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "Specified"
  environments                = [octopusdeploy_environment.env_1.id, octopusdeploy_environment.env_2.id]
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "tenanted_runbook_process" {
  runbook_id = octopusdeploy_runbook.tenanted_runbook.id
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
      run_on_server                      = true
    }
  }
}

resource "octopusdeploy_tag_set" "tagset_tag1" {
  name        = "tag1"
  description = "Test tagset"
  sort_order  = 0
}

resource "octopusdeploy_tag" "tag_a" {
  name        = "a"
  color       = "#333333"
  description = "tag a"
  sort_order  = 2
  tag_set_id  = octopusdeploy_tag_set.tagset_tag1.id
}

resource "octopusdeploy_tag" "tag_b" {
  name        = "b"
  color       = "#333333"
  description = "tag b"
  sort_order  = 3
  tag_set_id  = octopusdeploy_tag_set.tagset_tag1.id
}

resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
  tenant_tags = ["tag1/a", "tag1/b"]
  space_id    = var.octopus_space_id

  depends_on = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]

  project_environment {
    environments = [octopusdeploy_environment.env_1.id, octopusdeploy_environment.env_2.id]
    project_id   = octopusdeploy_project.tenanted.id
  }
}

resource "octopusdeploy_tenant" "tenant_team_b" {
  name        = "Team B"
  description = "Test tenant"
  tenant_tags = ["tag1/a", "tag1/b"]
  space_id    = var.octopus_space_id
  depends_on  = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]

  project_environment {
    environments = [octopusdeploy_environment.env_1.id, octopusdeploy_environment.env_2.id]
    project_id   = octopusdeploy_project.tenanted.id
  }
}

resource "octopusdeploy_project_scheduled_trigger" "tenanted_trigger" {
  name       = "Cron Tenanted"
  project_id = octopusdeploy_project.tenanted.id
  space_id   = octopusdeploy_project.tenanted.space_id
  tenant_ids = [octopusdeploy_tenant.tenant_team_a.id, octopusdeploy_tenant.tenant_team_b.id]
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "tenanted_runbook_trigger" {
  name        = "Cron Runbook"
  description = "This is a Cron schedule"
  project_id  = octopusdeploy_project.tenanted.id
  space_id    = octopusdeploy_project.tenanted.space_id
  tenant_ids  = [octopusdeploy_tenant.tenant_team_a.id, octopusdeploy_tenant.tenant_team_b.id]

  run_runbook_action {
    target_environment_ids = [octopusdeploy_environment.env_1.id]
    runbook_id             = octopusdeploy_runbook.tenanted_runbook.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}