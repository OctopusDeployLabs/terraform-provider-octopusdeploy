# basic deployment process with 2 run a script steps
resource "octopusdeploy_deployment_process" "example" {
  project_id = "Projects-123"
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
      run_on_server                      = true
    }
  }
  step {
    condition           = "Success"
    name                = "Hello world (using Bash)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartWithPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using Bash)"
      script_body                        = <<-EOT
          echo 'Hello world, using Bash'
          #TODO: Experiment with steps of your own :)
          echo '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
      run_on_server                      = true
    }
  }
}

# basic deployment process with 2 run a script steps as child steps
resource "octopusdeploy_deployment_process" "child_step_example" {
  project_id = "Projects-123"
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    target_roles        = [ "hello-world" ]
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
    }
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using Bash)"
      script_body                        = <<-EOT
          echo 'Hello world, using Bash'
          #TODO: Experiment with steps of your own :)
          echo '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
    }
  }
}

# rolling deployment process with a step with 2 run a script steps as child steps deploying to 2 targets in parallel
resource "octopusdeploy_deployment_process" "child_step_rolling_deployment_example" {
  project_id = "Projects-123"
  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    target_roles        = [ "hello-world" ]
    window_size         = 2
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell)"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
          #TODO: Experiment with steps of your own :)
          Write-Host '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
    }
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using Bash)"
      script_body                        = <<-EOT
          echo 'Hello world, using Bash'
          #TODO: Experiment with steps of your own :)
          echo '[Learn more about the types of steps available in Octopus](https://g.octopushq.com/OnboardingAddStepsLearnMore)'
        EOT
    }
  }
}

