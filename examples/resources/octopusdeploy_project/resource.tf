resource "octopusdeploy_project" "example" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "The development project."
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = "Lifecycles-123"
  name                                 = "Development Project (OK to Delete)"
  project_group_id                     = "ProjectGroups-123"
  tenanted_deployment_participation    = "TenantedOrUntenanted"

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }

  jira_service_management_extension_settings {
    connection_id             = "133d7fe602514060a48bc42ee9870f99"
    is_enabled                = false
    service_desk_project_name = "Test Service Desk Project (OK to Delete)"
  }

  servicenow_extension_settings {
    connection_id                       = "989034685e2c48c4b06a29286c9ef5cc"
    is_enabled                          = false
    is_state_automatically_transitioned = false
    standard_change_template_name       = "Standard Change Template Name (OK to Delete)"
  }

  template {
    default_value = "example-default-value"
    help_text     = "example-help-test"
    label         = "example-label"
    name          = "example-template-value"
    display_settings = {
      "Octopus.ControlType" : "SingleLineText"
    }
  }
}
