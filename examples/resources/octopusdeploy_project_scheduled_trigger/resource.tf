resource "octopusdeploy_project_scheduled_trigger" "once_daily_example" {
  name        = "Once Daily example"
  description = "This is a once daily schedule"
  project_id  = "projects-123"
  space_id    = "spaces-123"
  deploy_new_release_action {
    destination_environment_id = "environments-123"
  }
  once_daily_schedule {
    start_time   = "2024-03-22T09:00:00"
    days_of_week = ["Tuesday", "Wednesday", "Monday"]
  }
}

resource "octopusdeploy_project_scheduled_trigger" "continuous_example" {
  name        = "Continuous"
  description = "This is a continuous daily schedule"
  project_id  = "projects-123"
  space_id    = "spaces-123"
  deploy_new_release_action {
    destination_environment_id = "environments-123"
  }
  continuous_daily_schedule {
    interval      = "OnceHourly"
    hour_interval = 3
    run_after     = "2024-03-22T09:00:00"
    run_until     = "2024-03-29T13:00:00"
    days_of_week  = ["Monday", "Tuesday", "Friday"]
  }
}

resource "octopusdeploy_project_scheduled_trigger" "deploy_latest_example" {
  name       = "Deploy Latest"
  project_id  = "projects-123"
  space_id    = "spaces-123"
  deploy_latest_release_action {
    source_environment_id      = "environments-321"
    destination_environment_id = "environments-123"
    should_redeploy            = true
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "deploy_new_example" {
  name       = "Deploy New"
  project_id  = "projects-123"
  space_id    = "spaces-123"
  deploy_new_release_action {
    destination_environment_id = "environments-123"
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}

resource "octopusdeploy_project_scheduled_trigger" "runbook_example" {
  name        = "Runbook"
  description = "This is a Cron schedule"
  project_id  = "projects-123"
  space_id    = "spaces-123"
  run_runbook_action {
    target_environment_ids = ["environments-123", "environments-321"]
    runbook_id             = "runbooks-123"
  }
  cron_expression_schedule {
    cron_expression = "0 0 06 * * Mon-Fri"
  }
}