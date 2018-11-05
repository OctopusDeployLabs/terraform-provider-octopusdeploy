terraform-provider-octopusdeploy
---
A Terraform provider for [Octopus Deploy](https://octopus.com).

Based on the [go-octopusdeploy](https://github.com/MattHodge/go-octopusdeploy) Octopus Deploy client.

[![Build status](https://ci.appveyor.com/api/projects/status/a5ejcududsoug94e/branch/master?svg=true)](https://ci.appveyor.com/project/MattHodge/terraform-provider-octopusdeploy/branch/master)

> :warning: This provider is in heavy development. It is not production ready yet.

<!-- TOC -->

- [Go Dependencies](#go-dependencies)
- [Downloading & Installing](#downloading--installing)
- [Configure the Provider](#configure-the-provider)
- [Data Sources](#data-sources)
- [Provider Resources](#provider-resources)
- [Provider Resources (To Be Moved To /docs)](#provider-resources-to-be-moved-to-docs)
    - [Project Groups](#project-groups)
        - [Example Usage](#example-usage)
        - [Argument Reference](#argument-reference)
        - [Attributes Reference](#attributes-reference)
    - [Project](#project)
        - [Example Usage](#example-usage-1)
        - [Argument Reference](#argument-reference-1)
            - [Common Deployment Step Arguments](#common-deployment-step-arguments)
                - [Common Across All Deployment Steps](#common-across-all-deployment-steps)
                - [Configuration and Transformation](#configuration-and-transformation)
                - [Feed and Packages](#feed-and-packages)
                - [IIS Application Pool](#iis-application-pool)
        - [Attributes Reference](#attributes-reference-1)
    - [Variables](#variables)
        - [Example Usage](#example-usage-2)
        - [Argument Reference](#argument-reference-2)
        - [Attributes reference](#attributes-reference)
    - [Machine Policies](#machine-policies)
        - [Example Usage](#example-usage-3)
        - [Argument Reference](#argument-reference-3)
        - [Attributes Reference](#attributes-reference-2)
    - [Machines (Deployment Targets)](#machines-deployment-targets)
        - [Example Usage](#example-usage-4)
        - [Resource Argument Reference](#resource-argument-reference)
        - [Resource Attribute Reference](#resource-attribute-reference)
        - [Data Argument Reference](#data-argument-reference)
        - [Resource Attribute Reference](#resource-attribute-reference-1)

<!-- /TOC -->

# Go Dependencies
* Dependencies are managed using [dep](https://golang.github.io/dep/docs/new-project.html)

```bash
# Vendor new modules
dep ensure
```

# Downloading & Installing

As this provider is still under development, you will need to manually download it.

There are compiled binaries for most platforms in [Releases](https://github.com/MattHodge/terraform-provider-octopusdeploy/releases).

To use it, extract the binary for your platform into the same folder as your `.tf` file(s) will be located, then run `terraform init`.

# Configure the Provider

```hcl
# main.tf

provider "octopusdeploy" {
  address = "http://octopus.production.yolo"
  apikey  = "API-XXXXXXXXXXXXX"
}
```

# Data Sources

* [octopusdeploy_environment](docs/provider/data_sources/environment.md)

# Provider Resources

* [octopusdeploy_environment](docs/provider/resources/environment.md)


# Provider Resources (To Be Moved To /docs)
## Project Groups

[Project groups](https://octopus.com/docs/deployment-process/projects#project-group) are a way of organizing your projects.

### Example Usage

Basic usage:
```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}
```

Data usage:

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}
```

Basic usage with ID export used to create project:
```hcl
resource "octopusdeploy_project_group" "finance" {
  description = "Financial Applications"
  name        = "Finance"
}

resource "octopusdeploy_project" "billing_service" {
  description      = "The Finance teams billing service"
  lifecycle_id     = "Lifecycles-1"
  name             = "Billing Service"
  project_group_id = "${octopusdeploy_project_group.finance.id}"
}
```

### Argument Reference
* `description` - (Optional) Description of the project group
* `name` - (Required) Name of the project group

### Attributes Reference
* `id` - The ID of the project group


## Project

[Projects](https://octopus.com/docs/deployment-process/projects) can consist of multiple deployment steps, or they might only have a single step.

### Example Usage
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

  deployment_step_windows_service {
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

  deployment_step_inline_script {
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

  deployment_step_iis_website {
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

  deployment_step_package_script {
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
```

Which creates the following:

![Octopus Deploy Multiple Deployment Steps](https://i.imgur.com/yWRFjrU.png)

### Argument Reference
The following arguments define the project:
* `name` - (Required) Name of the project group.
* `description` - (Optional) Description of the project group.
* `lifecycle_id` - (Required) The ID of the lifecycle the project will use.
* `project_group_id` - (Required) The ID of the project group the project will be in.
* `default_failure_mode` - (Optional - Default is `EnvironmentDefault`) [Guided failure mode](https://octopus.com/docs/deployment-process/releases/guided-failures) tells Octopus that if something goes wrong during the deployment, instead of failing immediately, Octopus should ask for a human to intervene. Allowed values `EnvironmentDefault`, `Off`, `On`.
* `skip_machine_behavior` - (Optional - Default is `None`) Choose to skip or not skip deployment targets if they are unavailable during a deployment. Allowed values `SkipUnavailableMachines`, `None`.
* `deployment_step_windows_service` - (Optional) Creates a Windows Service deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `deployment_step_iis_website` - (Optional) Creates an IIS deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `deployment_step_inline_script` - (Optional) Creates inline script deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.
* `deployment_step_package_script` - (Optional) Creates package script deployment step. Can be specified multiple times in a project. Each block supports the fields documented below.

The `deployment_step_windows_service` block supports:
* `executable_path` - (Required) Path to the executable for the service
* `service_account` - (Optional - Default is `LocalSystem`) The account to run the service under
* `service_name` - (Required) The name of the service
* `service_start_mode` - (Optional - Default is `auto`) The start type for the service. Allowed values `auto`, `delayed-auto`, `demand`, `unchanged`
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section
* The arguments in the [Feed and Packages](#Feed-and-Packages) section
* The arguments in the [Configuration and Transformation](#Configuration-and-Transformation) section

The `deployment_step_iis_website` block supports:
* `anonymous_authentication` - (Optional - Default is `false`) Whether IIS should allow anonymous authentication.
* `basic_authentication` - (Optional - Default is `false`) Whether IIS should allow basic authentication with a 401 challenge.
* `website_name` - (Required) The name of the Website to be created.
* `windows_authentication` - (Optional - Default is `true`) Whether IIS should allow integrated Windows authentication with a 401 challenge.
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.
* The arguments in the [Feed and Packages](#Feed-and-Packages) section.
* The arguments in the [Configuration and Transformation](#Configuration-and-Transformation) section.
* The arguments in the [IIS Application Pool](#IIS-Application-Pool) section.

The `deployment_step_inline_script` block supports:
* `script_type` - (Required) The scripting language of the deployment step. Allowed values `PowerShell`, `CSharp`, `Bash`, `FSharp`.
* `script_body` - (Required) The script body. Multi-line strings are [supported by Terraform](https://www.terraform.io/docs/configuration/variables.html#strings).
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.

The `deployment_step_package_script` block supports:
* `script_file_name` - (Required) The script file name in the package.
* `script_parameters` - (Optional) Parameters expected by the script.
* The arguments in the [Common Across All Deployment Steps](#Common-Across-All-Deployment-Steps) section.
* The arguments in the [Feed and Packages](#Feed-and-Packages) section.

#### Common Deployment Step Arguments
The following arguments are shared amongst the `deployment_step` resources.
##### Common Across All Deployment Steps
* `step_condition` - (Optional - Default is `success`) Limit when this step will run by setting this condition.
* `step_name` - (Required) The name of the deployment step.
* `step_start_trigger` - (Optional - Default is `StartAfterPrevious`) Control whether the step waits for the previous step to complete, or runs parallel with it. Allowed values `StartAfterPrevious`, `StartWithPrevious`
* `target_roles` - (Required) A list of roles this deployment step will run on.
##### Configuration and Transformation
* `configuration_transforms` - (Optional - Default is `true`) Enables XML configuration transformations.
* `configuration_variables` - (Optional - Default is `true`) Enables replacing appSettings and connectionString entries in any .config file.
* `json_file_variable_replacement` - (Optional) A comma-separated list of file names to replace settings in, relative to the package contents.
##### Feed and Packages
* `feed_id` - (Optional - Default is `feeds-builtin`) The ID of the feed a package will be found in.
* `package` - (Required) ID / Name of the package to be deployed.
##### IIS Application Pool
* `application_pool_name` - (Required) Name of the application pool in IIS to create or reconfigure.
* `application_pool_framework` - (Optional - Default is `v4.0`) The version of the .NET common language runtime that this application. pool will use. Choose `v2.0` for applications built against .NET 2.0, 3.0 or 3.5. Choose `v4.0` for .NET 4.0 or 4.5.
* `application_pool_identity` - (Optional - Default is `ApplicationPoolIdentity`) Which built-in account will the application pool run under.

### Attributes Reference
* `deployment_process_id` - The ID of the projects deployment process.

## Variables

[Variables](https://octopus.com/docs/deployment-process/variables) are values that change based on the
scope of the deployments (e.g. changing SQL Connection Strings between production and staging deployments).

### Example Usage

Basic usage:

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

resource "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  type       = "String"
  value      = "Server=myServerAddress;Database=myDataBase;Trusted_Connection=True;"
}
```

More complex example (with environments and prompts)

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

data "octopusdeploy_environment" "staging" {
    name = "Staging"
}

resource "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  type       = "String"
  value      = "Server=myServerAddress;Database=myDataBase;Trusted_Connection=True;"

  scope {
    environments = ["${data.octopusdeploy_environment.staging.id}"]
  }

  prompt{
    label    = "SQL Server Connection String"
    required = false
  }
}
```

Data usage (with scope):

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

data "octopusdeploy_environment" "staging" {
    name = "Staging"
}

data "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  scope {
    environments = ["${data.octopusdeploy_environment.staging.id}"]
  }
}
```

### Argument Reference

* `project_id` (Required) ID of the Project to assign the variable against.
* `name` - (Required) Name of the variable
* `type` - (Required) Type of the variable. Must be one of `String`, `Certificate` or `AmazonWebServicesAccount` (`Sensitive` is not currently supported)
* `value` - (Required) The value of the variable
* `description` - (Optional) Description of the variable
* `scope` - (Optional) The scope to apply to this variable. Contains a list of arrays. All are optional:
    * (Optional) `environments`, `machines`, `actions`, `roles`, `channels`, `tenant_tags`
* `prompt` - (Optional) Prompt for value when a build is run
    * `label` - (Optional) The label for the prompt
    * `description` - (Optional) The description for the prompt
    * `required`- (Optional) Whether or not the value is required

### Attributes reference

* `id` - ID of the variable
* `name` - Name of the variable
* `type` - Type of the variable
* `value` - Value of the variable
* `description` - Description of the variable

## Machine Policies

[Machine policies](https://octopus.com/docs/infrastructure/machine-policies) are groups of settings that can be applied to Tentacle and SSH endpoints to modify their behavior.

Currently the Octopus terraform provider only provides machine policies as a data provider, as there are places elsewhere in
the provider that the IDs of machine policies need to be referenced.

### Example Usage

```hcl
data "octopusdeploy_machinepolicy" "default" {
  name = "Default Machine Policy"
}
```

### Argument Reference

* `name` - (Required) The name of the machine policy

### Attributes Reference

* `name` - The name of the machine policy
* `description` - The description of the machine policy
* `isdefault` - Whether or not this machine policy is the default policy

## Machines (Deployment Targets)

Octopus Deploy refers to Machines as [Deployment Targets](https://octopus.com/docs/infrastructure), however the API (and thus Terraform) refers to them as Machines.

### Example Usage

Basic Usage

```hcl
data "octopusdeploy_environment" "staging" {
  name = "Staging"
}

data "octopusdeploy_machinepolicy" "default" {
  name = "Default Machine Policy"
}

resource "octopusdeploy_machine" "testmachine" {
  name                            = "finance-web-01"
  environments                    = ["${data.octopusdeploy_environment.staging.id}"]
  isdisabled                      = false
  machinepolicy                   = "${data.octopusdeploy_machinepolicy.default.id}"
  roles                           = ["Staging"]
  tenanteddeploymentparticipation = "Untenanted"

  endpoint {
    communicationstyle = "TentaclePassive"
    thumbprint         = "81D0FF8B76FC"
    uri                = "https://finance-web-01:10933"
  }
}
```

Data Usage

```hcl
resource "octopusdeploy_machine" "testmachine" {
  name = "finance-web-01"
}
```

### Resource Argument Reference

* `name` - (Required) The name of the machine
* `endpoint` - (Required) The configuration of the machine endpoint
    * `communicationstyle` - (Required) Must be one of `None`, `TentaclePassive`, `TentacleActive`, `Ssh`, `OfflineDrop`, `AzureWebApp`, `Ftp`, `AzureCloudService`
    * `proxyid` - (Optional) ID of a defined proxy to use for communication with this machine
    * `thumbprint` - (Required) Thumbprint of the certificate this machine uses (if `communicationstyle` is `None` this should be blank)
    * `uri` - (Required) URI to access this machine (if `endpoint` is `None` this should be blank)
* `environments` - (Required) List of environment IDs to be assigned to this machine
* `isdisabled` - (Required) Whether or not this machine is disabled
* `machinepolicy` - (Required) The ID of the machine policy to be assigned to this machine
* `roles` - (Required) List of the roles to be assigned to this machine
* `tenanteddeploymentparticipation` - (Required) Must be one of `Untenanted`, `TenantedOrUntenanted`, `Tenanted`
* `tenantids` - (Optional) If tenanted, a list of the tenant IDs for this machine
* `tenanttags` - (Optional) If tenanted, a list of the tenant tags for this machine

### Resource Attribute Reference

* `environments` - List of environment IDs this machine is assigned to
* `haslatestcalamari` - Whether or not this machine has the latest Calamari version
* `isdisabled` - Whether or not this machine is disabled
* `isinprocess` - Whether or not this machine is being processed
* `machinepolicy` - The ID of the machine policy assigned to this machine
* `roles` - A list of the roles assigned to this machine
* `status` - The machine status code
* `statussummary` - Plain text description of the machine status
* `tenanteddeploymentparticipation` - One of `Untenanted`, `TenantedOrUntenanted`, `Tenanted`
* `tenantids` - If tenanted, a list of the tenant IDs for this machine
* `tenanttags` -  If tenanted, a list of the tenant tags for this machine

### Data Argument Reference

* `name` - (Required) The name of the machine

### Resource Attribute Reference

All items from the Resource Attribute Reference, and additionally:

* `endpoint_communicationstyle` - One of `None`, `TentaclePassive`, `TentacleActive`, `Ssh`, `OfflineDrop`, `AzureWebApp`, `Ftp`, `AzureCloudService`
* `endpoint_proxyid` - ID of a defined proxy to use for communication with this machine
* `endpoint_tentacleversiondetails_upgradelocked` - Whether or not this machine tentacle is upgrade locked
* `endpoint_tentacleversiondetails_upgraderequired` - Whether or not this machine tentacle required an upgrade
* `endpoint_tentacleversiondetails_upgradesuggested` - Whether or not this machine tentacle has a suggested ugrade
* `endpoint_tentacleversiondetails_version` - Version number of this machine tentacle
* `endpoint_thumbprint` - Thumbprint of the certificate this machine uses
* `endpoint_uri` - URI to access this machine
