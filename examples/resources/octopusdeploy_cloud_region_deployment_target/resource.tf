resource "octopusdeploy_cloud_region_deployment_target" "example" {
  environments = ["Environments-123", "Environment-321"]
  name         = "Azure Web App Deployment Target (OK to Delete)"
  roles        = ["Development Team", "System Administrators"]
}
