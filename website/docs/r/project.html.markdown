---
layout: "octopusdeploy"
page_title: "Octopus Deploy: project"
---

# Resource: project

[Projects](https://octopus.com/docs/deployment-process/projects) can consist of multiple deployment steps, or they might only have a single step.

## Example Usage

Basic project with some settings but no deployment steps:

```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}

resource "octopusdeploy_project" "test_project" {
  description           = "A really groundbreaking app"
  lifecycle_id          = "Lifecycles-1"
  name                  = "Epic App"
  project_group_id      = "${octopusdeploy_project_group.finance.id}"
  skip_machine_behavior = "SkipUnavailableMachines"
}
```

Project with many settings configured and multiple types of deployment steps:

```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}

resource "octopusdeploy_project" "billing_service" {
  description           = "Billing Frontend and Backend"
  lifecycle_id          = "Lifecycles-1"
  name                  = "Billing Service"
  project_group_id      = "${octopusdeploy_project_group.finance.id}"
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

```

Which creates the following:

![Octopus Deploy Multiple Deployment Steps](https://i.imgur.com/yWRFjrU.png)

## Argument Reference

The following arguments define the project:

* `name` - (Required) Name of the project group.
* `description` - (Optional) Description of the project group.
* `lifecycle_id` - (Required) The ID of the lifecycle the project will use.
* `project_group_id` - (Required) The ID of the project group the project will be in.
* `default_failure_mode` - (Optional - Default is `EnvironmentDefault`) [Guided failure mode](https://octopus.com/docs/deployment-process/releases/guided-failures) tells Octopus that if something goes wrong during the deployment, instead of failing immediately, Octopus should ask for a human to intervene. Allowed values `EnvironmentDefault`, `Off`, `On`.
* `skip_machine_behavior` - (Optional - Default is `None`) Choose to skip or not skip deployment targets if they are unavailable during a deployment. Allowed values `SkipUnavailableMachines`, `None`.
* `deployment_step` - (Optional) Defines the process deployment step(s).
* `windows_service` - (Optional) Creates a Windows Service deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `iis_website` - (Optional) Creates an IIS deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `inline_script` - (Optional) Creates inline script deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `package_script` - (Optional) Creates package script deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.

The `windows_service` block supports:

* `executable_path` - (Required) Path to the executable for the service
* `service_account` - (Optional - Default is `LocalSystem`) The account to run the service under
* `service_name` - (Required) The name of the service
* `service_start_mode` - (Optional - Default is `auto`) The start type for the service. Allowed values `auto`, `delayed-auto`, `demand`, `unchanged`
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section
* The arguments in the [Feed and Packages](#Feed-and-Packages) section
* The arguments in the [Configuration and Transformation](#Configuration-and-Transformation) section

The `iis_website` block supports:

* `anonymous_authentication` - (Optional - Default is `false`) Whether IIS should allow anonymous authentication.
* `basic_authentication` - (Optional - Default is `false`) Whether IIS should allow basic authentication with a 401 challenge.
* `website_name` - (Required) The name of the Website to be created.
* `windows_authentication` - (Optional - Default is `true`) Whether IIS should allow integrated Windows authentication with a 401 challenge.
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.
* The arguments in the [Feed and Packages](#Feed-and-Packages) section.
* The arguments in the [Configuration and Transformation](#Configuration-and-Transformation) section.
* The arguments in the [IIS Application Pool](#IIS-Application-Pool) section.

The `inline_script` block supports:

* `script_type` - (Required) The scripting language of the deployment step. Allowed values `PowerShell`, `CSharp`, `Bash`, `FSharp`.
* `script_body` - (Required) The script body. Multi-line strings are [supported by Terraform](https://www.terraform.io/docs/configuration/variables.html#strings).
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.

The `package_script` block supports:

* `script_file_name` - (Required) The script file name in the package.
* `script_parameters` - (Optional) Parameters expected by the script.
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.
* The arguments in the [Feed and Packages](#Feed-and-Packages) section.

### Common Deployment Step Arguments

The following arguments are shared amongst the `deployment_step` resources.

#### Common Across All Deployment Steps

* `step_condition` - (Optional - Default is `success`) Limit when this step will run by setting this condition.
* `step_name` - (Required) The name of the deployment step.
* `step_start_trigger` - (Optional - Default is `StartAfterPrevious`) Control whether the step waits for the previous step to complete, or runs parallel with it. Allowed values `StartAfterPrevious`, `StartWithPrevious`
* `target_roles` - (Required) A list of roles this deployment step will run on.

#### Configuration and Transformation

* `configuration_transforms` - (Optional - Default is `true`) Enables XML configuration transformations.
* `configuration_variables` - (Optional - Default is `true`) Enables replacing appSettings and connectionString entries in any .config file.
* `json_file_variable_replacement` - (Optional) A comma-separated list of file names to replace settings in, relative to the package contents.

#### Feed and Packages

* `feed_id` - (Optional - Default is `feeds-builtin`) The ID of the feed a package will be found in.
* `package` - (Required) ID / Name of the package to be deployed.

#### IIS Application Pool

* `application_pool_name` - (Required) Name of the application pool in IIS to create or reconfigure.
* `application_pool_framework` - (Optional - Default is `v4.0`) The version of the .NET common language runtime that this application. pool will use. Choose `v2.0` for applications built against .NET 2.0, 3.0 or 3.5. Choose `v4.0` for .NET 4.0 or 4.5.
* `application_pool_identity` - (Optional - Default is `ApplicationPoolIdentity`) Which built-in account will the application pool run under.

### Attributes Reference

* `deployment_process_id` - The ID of the projects deployment process.
