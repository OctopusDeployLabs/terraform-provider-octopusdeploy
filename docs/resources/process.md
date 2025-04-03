---
page_title: "octopusdeploy_process Resource - terraform-provider-octopusdeploy"
subcategory: "Runbook & Deployment Processes"
description: |-
  This resource manages Runbook and Deployment Processes in Octopus Deploy. It's used in collaboration with octopusdeploy_process_step and octopusdeploy_process_step_order.
---

# octopusdeploy_process (Resource)

This resource manages Runbook and Deployment Processes in Octopus Deploy. It's used in collaboration with `octopusdeploy_process_step` and `octopusdeploy_process_step_order`.

~> This resource is the successor to the original `octopusdeploy_deployment_process` resource, which suffered from numerous problems including: state drift, data inconsistency when reordering or inserting steps, and lack of awareness of Version-Controlled projects.

### Remarks

The `octopusdeploy_process` resource is used in conjunction with a series of other building-block resources to form a full process. They are deliberately designed with dependencies between them so that the deployment process will be incrementally "built up" in Octopus with a series of incremental updates. You can use only the building blocks you need for your process (i.e. if your process doesn't involve Child Steps, you don't need to deal with the `octopusdeploy_process_child_step` resource.)

At a minimum, to get a functional deployment process, you will need:

1. A Project (`octopusdeploy_project`)
1. A Deployment Process referencing the Project (`octopusdeploy_process`)
1. One or more Steps referencing the Deployment Process (`octopusdeploy_process_step`)

The `octopusdeploy_process_step_order` resource isn't strictly required, but it's highly recommended. If you need to change the order of steps in your process, or insert a new step within an existing process, you'll need the Step Order defined first.

Without a defined Step Order, the Steps will be added to the process in the order they're applied by Terraform. This is usually the order they appear in your HCL, but may not always be deterministic. 

## Example Usage
~> See the docs for `octopusdeploy_process_step`, `octopusdeploy_process_steps_order`, `octopusdeploy_process_child_step` and `octopusdeploy_process_child_steps_order` for more detailed examples.

### Deployment Process
```terraform
# Example of a Deployment Process with three steps and an explicit Step Order
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
}

resource "octopusdeploy_process_step" "run_script" {
  # Run script step
  process_id  = octopusdeploy_process.example.id  
  name = "Run My Script"
  properties = {
    "Octopus.Action.MaxParallelism" = "2"
    "Octopus.Action.TargetRoles" = "role-1,role-2"
  }
  type = "Octopus.Script"
  environments = [octopusdeploy_environment.development.id]
  excluded_environments = [octopusdeploy_environment.production.id]
  channels = [octopusdeploy_channel.example.id]
  notes = "Script example"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      Write-Host "Executing step..."
    EOT
  }
}

resource "octopusdeploy_process_step" "approval" {
  # Manual intervention
  process_id  = octopusdeploy_process.example.id
  name = "Approve deployment"
  type = "Octopus.Manual"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Manual.Instructions" = "Example of manual blocking step"
    "Octopus.Action.Manual.BlockConcurrentDeployments" = "True"
    "Octopus.Action.Manual.ResponsibleTeamIds" = "teams-managers"
  }
}

resource "octopusdeploy_process_step" "deploy_package" {
  # Package deployment with primary package
  process_id  = octopusdeploy_process.example.id
  name = "Package deployment"
  properties = {
    "Octopus.Action.TargetRoles" = "role-one"
  }
  type = "Octopus.TentaclePackage"
  packages = {
    "": {
      package_id: "my.package"
      feed_id: "Feeds-1"
    }
  }
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"    
    # Reference primary package in execution properties for legacy purposes
    "Octopus.Action.Package.DownloadOnTentacle" = "False"
    "Octopus.Action.Package.FeedId" = "Feeds-1"
    "Octopus.Action.Package.PackageId" = "my.package"
  }
}

resource "octopusdeploy_process_steps_order" "example" {
  process_id  = octopusdeploy_process.example.id
  steps = [
    octopusdeploy_process_step.run_script.id,
    octopusdeploy_process_step.approval.id,
    octopusdeploy_process_step.deploy_package.id,
  ]
}
```

### Runbook Process
```terraform
# Example of a Runbook Process with two steps and an explicit Step Order
# To manage a Runbook process, specify both the Project and Runbook IDs (usually via Terraform resource references)
resource "octopusdeploy_process" "example" {
  space_id = "Spaces-1"
  project_id  = "Projects-21"
  runbook_id  = "Runbooks-42"
}

resource "octopusdeploy_process_step" "run_script" {
  # Run script step
  process_id  = octopusdeploy_process.example.id  
  name = "Run My Script"
  properties = {
    "Octopus.Action.MaxParallelism" = "2"
    "Octopus.Action.TargetRoles" = "role-1,role-2"
  }
  type = "Octopus.Script"
  environments = [octopusdeploy_environment.development.id]
  excluded_environments = [octopusdeploy_environment.production.id]
  channels = [octopusdeploy_channel.example.id]
  notes = "Script example"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      Write-Host "Executing step..."
    EOT
  }
}

resource "octopusdeploy_process_step" "deploy_package" {
  # Package deployment with primary package
  process_id  = octopusdeploy_process.example.id
  name = "Package deployment"
  properties = {
    "Octopus.Action.TargetRoles" = "role-one"
  }
  type = "Octopus.TentaclePackage"
  packages = {
    "": {
      package_id: "my.package"
      feed_id: "Feeds-1"
    }
  }
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"    
    # Reference primary package in execution properties for legacy purposes
    "Octopus.Action.Package.DownloadOnTentacle" = "False"
    "Octopus.Action.Package.FeedId" = "Feeds-1"
    "Octopus.Action.Package.PackageId" = "my.package"
  }
}

resource "octopusdeploy_process_steps_order" "example" {
  process_id  = octopusdeploy_process.example.id
  steps = [
    octopusdeploy_process_step.run_script.id,
    octopusdeploy_process_step.deploy_package.id,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Id of the project this process belongs to.

### Optional

- `runbook_id` (String) Id of the runbook this process belongs to. When not set this resource represents deployment process of the project
- `space_id` (String) The space ID associated with this process.

### Read-Only

- `id` (String) The unique ID for this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import [options] octopusdeploy_process.<name> <process-id>
```
