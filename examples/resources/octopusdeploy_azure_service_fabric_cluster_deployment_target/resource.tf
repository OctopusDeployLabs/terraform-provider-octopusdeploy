resource "octopusdeploy_azure_service_fabric_cluster_deployment_target" "example" {
  connection_endpoint               = "[connection-endpoint]"
  environments                      = ["Environments-123", "Environment-321"]
  name                              = "Azure Service Fabric Cluster Deployment Target (OK to Delete)"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
}
