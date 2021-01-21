resource "octopusdeploy_azure_web_app_deployment_target" "example" {
  account_id                        = "Accounts-123"
  name                              = "Azure Web App Deployment Target (OK to Delete)"
  resource_group_name               = "[resource-group-name]"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
  web_app_name                      = "[web_app_name]"
}