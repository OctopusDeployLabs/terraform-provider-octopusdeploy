resource "octopusdeploy_project_group" "example" {
  name        = "Example"
  description = "Example Group"
}

resource "octopusdeploy_project" "example" {
  name                                 = "Example"
  lifecycle_id                         = "Lifecycles-101"
  project_group_id                     = octopusdeploy_project_group.example.id
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Project with Built-In Trigger"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  tenanted_deployment_participation    = "Untenanted"
  included_library_variable_sets       = []

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_channel" "example" {
    name = "Example Channel"
    project_id = octopusdeploy_project.example.id
    lifecycle_id = "Lifecycles-101"
}

data "octopusdeploy_feeds" "built_in" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
}

resource "octopusdeploy_deployment_process" "example" {
  project_id = octopusdeploy_project.example.id
  step {
    condition           = "Success"
    name                = "Step One"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"
    run_script_action {
      condition                          = "Success"
      is_disabled                        = false
      is_required                        = true
      name                               = "Action One"
      script_body                        = <<-EOT
          $ExtractedPath = $OctopusParameters["Octopus.Action.Package[my.package].ExtractedPath"]
          Write-Host $ExtractedPath
        EOT
      run_on_server                      = true

      package {
        name                      = "my.package"
        package_id                = "my.package"
        feed_id                   = data.octopusdeploy_feeds.built_in.feeds[0].id
        acquisition_location      = "Server"
        extract_during_deployment = true
      }
    }
  }
}

resource "octopusdeploy_built_in_trigger" "example" {
  project_id = octopusdeploy_project.example.id
  channel_id = octopusdeploy_channel.example.id
  
  release_creation_package = {
    deployment_action = "Action One"
    package_reference = "my.package"
  }

  depends_on = [
    octopusdeploy_project.example,
    octopusdeploy_channel.example,
    octopusdeploy_deployment_process.example
  ]
}
