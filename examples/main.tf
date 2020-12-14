/*provider "octopusdeploy" {
  address = "http://localhost:8081/"
  apikey  = "API-XXXXXXXXXXXXXXXXXXXX"
}*/

resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}

resource "octopusdeploy_project_deployment_target_trigger" "auto_deploy_trigger" {
  name             = "Auto Deploy"
  project_id       = octopusdeploy_project.billing_service.id
  event_groups     = ["Machine"]
  event_categories = ["MachineCleanupFailed"]
  should_redeploy  = true

  roles = [
    "Billing-Batch-Processor",
  ]
}

resource "octopusdeploy_project" "billing_service" {
  description           = "Billing Frontend and Backend"
  lifecycle_id          = "Lifecycles-1"
  name                  = "Billing Service"
  project_group_id      = octopusdeploy_project_group.finance.id
  skip_machine_behavior = "SkipUnavailableMachines"

  deployment_step {
    windows_service {
      executable_path                = "batch_processor\\batch_processor_service.exe"
      service_name                   = "Billing Batch Processor"
      step_name                      = "Deploy Billing Batch Processor Windows Service"
      step_condition                 = "failure"
      package                        = "Billing.BatchProcessor"
      json_file_variable_replacement = "appsettings.json"

      target_roles = [
        "Billing-Batch-Processor",
      ]
    }

    inline_script {
      step_name   = "Cleanup Temporary Files"
      script_type = "PowerShell"

      script_body = <<EOF
$oldFiles = Get-ChildItem -Path 'C:\billing_archived_jobs'
Remove-Item $oldFiles -Force -Recurse
EOF

      target_roles = [
        "Billing-Batch-Processor",
      ]
    }

    iis_website {
      step_name                  = "Deploy Billing API"
      website_name               = "Billing API"
      application_pool_name      = "Billing"
      application_pool_framework = "v2.0"
      basic_authentication       = true
      windows_authentication     = false
      package                    = "Billing.API"

      target_roles = [
        "Billing-API-Asia",
        "Billing-API-Europe",
      ]
    }

    package_script {
      step_name         = "Verify API Deployment"
      package           = "Billing.API"
      script_file_name  = "scripts\\verify_deployment.ps1"
      script_parameters = "-Verbose"

      target_roles = [
        "Billing-API-Asia",
        "Billing-API-Europe",
      ]
    }
  }
}
