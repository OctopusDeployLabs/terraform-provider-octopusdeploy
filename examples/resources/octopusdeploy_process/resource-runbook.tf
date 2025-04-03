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
