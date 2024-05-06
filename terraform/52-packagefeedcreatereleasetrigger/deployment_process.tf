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
  step {
    condition           = "Success"
    name                = "Helm one"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.HelmChartUpgrade"
      name                               = "Helm one"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      worker_pool_id                     = ""
      properties = {
        "Octopus.Action.Helm.ClientVersion"         = "V3"
        "Octopus.Action.Helm.Namespace"             = "dev"
        "Octopus.Action.Package.DownloadOnTentacle" = "False"
        "Octopus.Action.Helm.ResetValues"           = "True"
      }
      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []

      primary_package {
        package_id           = "redis"
        acquisition_location = "Server"
        feed_id              = "${octopusdeploy_helm_feed.feed_helm_charts.id}"
        properties           = { SelectionMode = "immediate" }
      }

      features = []
    }

    properties   = {}
    target_roles = ["k8s"]
  }
  step {
    condition           = "Success"
    name                = "Helm two"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.HelmChartUpgrade"
      name                               = "Helm two"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      worker_pool_id                     = ""
      properties = {
        "Octopus.Action.Helm.ClientVersion"         = "V3"
        "Octopus.Action.Helm.Namespace"             = "dev"
        "Octopus.Action.Package.DownloadOnTentacle" = "False"
        "Octopus.Action.Helm.ResetValues"           = "True"
      }
      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []

      primary_package {
        package_id           = "prometheus"
        acquisition_location = "Server"
        feed_id              = "${octopusdeploy_helm_feed.feed_helm_charts.id}"
        properties           = { SelectionMode = "immediate" }
      }

      features = []
    }

    properties   = {}
    target_roles = ["k8s"]
  }
}
