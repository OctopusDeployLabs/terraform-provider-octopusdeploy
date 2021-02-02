resource "octopusdeploy_deployment_process" "example" {
  project_id = "Projects-123"
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    action {
      action_type                        = "Octopus.Script"
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      properties = {
        "Octopus.Action.RunOnServer"         = "true"
        "Octopus.Action.Script.ScriptBody"   = <<-EOT
                    Write-Host 'Hello world, using PowerShell'
                    #TODO: Experiment with steps of your own :)
                    Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
                EOT
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax"       = "PowerShell"
      }
      run_on_server = false
    }
  }
  step {
    condition           = "Success"
    name                = "Hello world (using Bash)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartWithPrevious"
    action {
      action_type                        = "Octopus.Script"
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using Bash)"
      properties = {
        "Octopus.Action.RunOnServer"         = "true"
        "Octopus.Action.Script.ScriptBody"   = <<-EOT
                    echo 'Hello world, using Bash'
                    #TODO: Experiment with steps of your own :)
                    echo '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
                EOT
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax"       = "Bash"
      }
      run_on_server = false
    }
  }
}
