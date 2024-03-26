data "octopusdeploy_feeds" "built_in_feed" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
}

resource "octopusdeploy_deployment_process" "example" {
  project_id = "${octopusdeploy_project.deploy_frontend_project.id}"
  step {
    condition           = "Success"
    name                = "Dummy step 1"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell) 1"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
        EOT
      run_on_server                      = true

      package {
        name                      = "nginx"
        package_id                = "nginx"
        feed_id                   = octopusdeploy_docker_container_registry.docker_feed.id
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
      package {
        name                      = "busybox"
        package_id                = "busybox"
        feed_id                   = octopusdeploy_docker_container_registry.docker_feed.id
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
    }
  }
  step {
    condition           = "Success"
    name                = "Dummy step 2"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      can_be_used_for_project_versioning = false
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Hello world (using PowerShell) 2"
      script_body                        = <<-EOT
          Write-Host 'Hello world, using PowerShell'
        EOT
      run_on_server                      = true

      package {
        name                      = "scratch"
        package_id                = "scratch"
        feed_id                   = octopusdeploy_docker_container_registry.docker_feed.id
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
    }
  }
}
