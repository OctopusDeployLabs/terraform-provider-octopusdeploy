resource "octopusdeploy_process" "example" {
  owner_id  = "Projects-12"
}

resource "octopusdeploy_process_step" "example" {
  process_id  = octopusdeploy_process.example.id
  name = "Example Step"
  type = "Octopus.Script"
  properties = {
    "Octopus.Action.TargetRoles" = "role-1"
  }
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Example Step...'"
  }
}

resource "octopusdeploy_process_child_step" "child_one" {
  process_id  = octopusdeploy_process.example.id
  parent_id = octopusdeploy_process_step.example.id
  name = "Child One"
  type = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Child 1...'"
  }
}

resource "octopusdeploy_process_child_step" "child_two" {
  process_id  = octopusdeploy_process.example.id
  parent_id = octopusdeploy_process_step.example.id
  name = "Child Two"
  type = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Child 3...'"
  }
}

resource "octopusdeploy_process_child_steps_order" "example" {
  process_id = octopusdeploy_process.example.id
  parent_id = octopusdeploy_process_step.example.id
  children = [
    octopusdeploy_process_child_step.child_one.id,
    octopusdeploy_process_child_step.child_two.id,
  ]
}

