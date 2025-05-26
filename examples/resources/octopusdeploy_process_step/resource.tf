resource "octopusdeploy_environment" "development" {
  name = "Development"
}

resource "octopusdeploy_environment" "production" {
  name = "Production"
}

resource "octopusdeploy_project" "example" {
  project_group_id = "ProjectGroups-1"
  lifecycle_id = "Lifecycles-1"
  name = "Example"
}

resource "octopusdeploy_channel" "example" {
  name       = "Example Channel (OK to Delete)"
  project_id = octopusdeploy_project.example.id
}

resource "octopusdeploy_process" "example" {
  project_id  = octopusdeploy_project.example.id
}

# Run script step
resource "octopusdeploy_process_step" "run_script" {
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

# Manual intervention
resource "octopusdeploy_process_step" "approval" {
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

# Package deployment with primary package
resource "octopusdeploy_process_step" "deploy_package" {
  process_id  = octopusdeploy_process.example.id
  name = "Package deployment"
  properties = {
    "Octopus.Action.TargetRoles" = "role-one"
  }
  type = "Octopus.TentaclePackage"
  primary_package = {
    package_id: "my.package"
    feed_id: "Feeds-1"
  }
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"    
  }
}
