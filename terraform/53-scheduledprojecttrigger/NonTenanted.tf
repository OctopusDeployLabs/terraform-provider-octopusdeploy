resource "octopusdeploy_project" "non_tenanted" {
  space_id                             = var.octopus_space_id
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Non Tenanted"
  project_group_id                     = octopusdeploy_project_group.project_group_test.id
  tenanted_deployment_participation    = "Untenanted"
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


resource "octopusdeploy_deployment_process" "non_tenanted_deployment_process" {
  project_id = octopusdeploy_project.non_tenanted.id
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

resource "octopusdeploy_runbook" "non_tenanted_runbook" {
  project_id         = octopusdeploy_project.non_tenanted.id
  name               = "Runbook"
  description        = "Test Runbook"
  multi_tenancy_mode = "Untenanted"
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

resource "octopusdeploy_runbook_process" "non_tenanted_runbook_process" {
  runbook_id = octopusdeploy_runbook.non_tenanted_runbook.id
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

resource "octopusdeploy_project_scheduled_trigger" "once_daily" {
  name        = "Once Daily"
  description = "This is a once daily schedule"
  project_id  = octopusdeploy_project.non_tenanted.id
  space_id    = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  once_daily_schedule {
    start_time   = "2024-03-22T09:00:00"
    days_of_week = ["Tuesday", "Wednesday", "Monday"]
  }
}

resource "octopusdeploy_project_scheduled_trigger" "continous" {
  name        = "Continuous"
  description = "This is a continuous daily schedule"
  project_id  = octopusdeploy_project.non_tenanted.id
  space_id    = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  continuous_daily_schedule {
    interval      = "OnceHourly"
    hour_interval = 3
    run_after     = "2024-03-22T09:00:00"
    run_until     = "2024-03-29T13:00:00"
    days_of_week  = ["Saturday", "Monday", "Sunday", "Tuesday", "Thursday", "Friday"]
  }
}

resource "octopusdeploy_project_scheduled_trigger" "days_per_month_date_of_month" {
  name       = "Days Per Month"
  project_id = octopusdeploy_project.non_tenanted.id
  space_id   = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  days_per_month_schedule {
    start_time            = "2024-03-22T09:00:00"
    monthly_schedule_type = "DateOfMonth"
    date_of_month         = "31"
  }
}


resource "octopusdeploy_project_scheduled_trigger" "days_per_month_day_of_month" {
  name       = "Days Per Month Specific Day"
  project_id = octopusdeploy_project.non_tenanted.id
  space_id   = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  days_per_month_schedule {
    start_time            = "2024-03-22T09:00:00"
    monthly_schedule_type = "DayOfMonth"
    day_number_of_month   = "L"
    day_of_week           = "Monday"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "cron" {
  name        = "Cron"
  description = "This is a Cron schedule"
  project_id  = octopusdeploy_project.non_tenanted.id
  space_id    = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "deploy_latest" {
  name       = "Cron Deploy Latest"
  project_id = octopusdeploy_project.non_tenanted.id
  space_id   = octopusdeploy_project.non_tenanted.space_id
  deploy_latest_release_action {
    source_environment_id      = octopusdeploy_environment.env_1.id
    destination_environment_id = octopusdeploy_environment.env_2.id
    should_redeploy            = true
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "deploy_new" {
  name       = "Cron Deploy New"
  project_id = octopusdeploy_project.non_tenanted.id
  space_id   = octopusdeploy_project.non_tenanted.space_id
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}


resource "octopusdeploy_project_scheduled_trigger" "timezone" {
  name       = "Specific Timezone"
  project_id = octopusdeploy_project.non_tenanted.id
  space_id   = octopusdeploy_project.non_tenanted.space_id
  timezone   = "Australia/Sydney"
  deploy_new_release_action {
    destination_environment_id = octopusdeploy_environment.env_1.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "runbook" {
  name        = "Cron Runbook"
  description = "This is a Cron schedule"
  project_id  = octopusdeploy_project.non_tenanted.id
  space_id    = octopusdeploy_project.non_tenanted.space_id
  run_runbook_action {
    target_environment_ids = [octopusdeploy_environment.env_1.id]
    runbook_id             = octopusdeploy_runbook.non_tenanted_runbook.id
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}
