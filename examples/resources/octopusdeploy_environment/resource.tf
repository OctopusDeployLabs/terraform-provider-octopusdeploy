resource "octopusdeploy_environment" "example" {
  allow_dynamic_infrastructure = false
  description                  = "An environment for the development team."
  name                         = "Development Environment (OK to Delete)"
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
}
