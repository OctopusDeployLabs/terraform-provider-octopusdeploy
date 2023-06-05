data "octopusdeploy_azure_service_fabric_cluster_deployment_targets" "data_lookup" {
  partial_name = "Service Fabric"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_azure_service_fabric_cluster_deployment_targets.data_lookup.azure_service_fabric_cluster_deployment_targets[0].id
}