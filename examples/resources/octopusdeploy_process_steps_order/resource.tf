resource "octopusdeploy_process" "example" {
  owner_id  = "Projects-12"
}

resource "octopusdeploy_process_step" "one" {
  process_id  = octopusdeploy_process.example.id
  name = "Step One"
  type = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Step 1...'"
  }
}

resource "octopusdeploy_process_step" "two" {
  process_id  = octopusdeploy_process.example.id
  name = "Step Two"
  type = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Step 2...'"
  }
}

resource "octopusdeploy_process_step" "three" {
  process_id  = octopusdeploy_process.example.id
  name = "Step Three"
  type = "Octopus.Script"
  execution_properties = {
    "Octopus.Action.Script.ScriptBody" = "Write-Host 'Step 3...'"
  }
}

resource "octopusdeploy_process_steps_order" "example" {
  process_id  = octopusdeploy_process.example.id
  steps = [
    octopusdeploy_process_step.one.id,
    octopusdeploy_process_step.two.id,
    octopusdeploy_process_step.three.id,
  ]
}
