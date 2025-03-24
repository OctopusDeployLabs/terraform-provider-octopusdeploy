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
    owner_id  = octopusdeploy_project.example.id
}

resource "octopusdeploy_process_step" "parent" {
  process_id  = octopusdeploy_process.example.id
  name = "Parent step"
  properties = {
    "Octopus.Action.TargetRoles" = "role-1"
  }
  type = "Octopus.Manual"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Manual.Instructions" = "Approve before executing child steps"
    "Octopus.Action.Manual.BlockConcurrentDeployments" = "True"
    "Octopus.Action.Manual.ResponsibleTeamIds" = "teams-managers"
  }
}

resource "octopusdeploy_process_child_step" "run_script" {
  process_id  = octopusdeploy_process.example.id
  parent_id = octopusdeploy_process_step.parent.id
  name = "Run My Script"
  type = "Octopus.Script"
  environments = [octopusdeploy_environment.development.id]
  excluded_environments = [octopusdeploy_environment.production.id]
  channels = [octopusdeploy_channel.example.id]
  notes = "Child script example"
  execution_properties = {
    "Octopus.Action.RunOnServer" = "True"
    "Octopus.Action.Script.ScriptSource" = "Inline"
    "Octopus.Action.Script.Syntax"       = "PowerShell"
    "Octopus.Action.Script.ScriptBody" = <<-EOT
      Write-Host "Executing child step after approval..."
    EOT
  }
}
