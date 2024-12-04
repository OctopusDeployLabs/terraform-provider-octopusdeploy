resource "octopusdeploy_project_group" "tp" {
  name        = "DevOps Projects"
  description = "My DevOps projects group"
}

resource "octopusdeploy_project" "tp" {
  name             = "My DevOps Project"
  description      = "test project"
  lifecycle_id     = "Lifecycles-1"
  project_group_id = octopusdeploy_project_group.tp.id

  depends_on  = [octopusdeploy_project_group.tp]
}

resource "octopusdeploy_deployment_process" "process" {
  project_id = octopusdeploy_project.tp.id

  step {
    name = "Hello World"
    target_roles        = [ "hello-world" ]
    start_trigger       = "StartAfterPrevious"
    package_requirement = "LetOctopusDecide"
    condition           = "Success"

    run_script_action {
      name                               = "Hello World"
      is_disabled                        = false
      is_required                        = true
      script_body                        = "Write-Host 'hello world'"
      script_syntax                      = "PowerShell"
      can_be_used_for_project_versioning = true
      sort_order = 1


      package {
        name                      = "Package"
        feed_id                   = "feeds-builtin"
        package_id                = "myExpressApp"
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
    }
  }

  depends_on  = [octopusdeploy_project.tp]
}

resource "octopusdeploy_project_versioning_strategy" "tp" {
  project_id = octopusdeploy_project.tp.id
  space_id = octopusdeploy_project.tp.space_id
  donor_package_step_id = octopusdeploy_deployment_process.process.step[0].run_script_action[0].id
  donor_package = {
    deployment_action = "Hello World"
    package_reference = "Package"
  }
  depends_on = [
    octopusdeploy_project_group.tp,
    octopusdeploy_deployment_process.process
  ]
}
