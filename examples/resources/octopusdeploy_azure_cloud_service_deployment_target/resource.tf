resource "octopusdeploy_azure_cloud_service_deployment_target" "example" {
  account_id                        = "Accounts-123"
  cloud_service_name                = "[cloud_service_name]"
  environments                      = ["Environments-123", "Environment-321"]
  name                              = "Azure Cloud Service Deployment Target (OK to Delete)"
  storage_account_name              = "[storage_account_name]"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
}