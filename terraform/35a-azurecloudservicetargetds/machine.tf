data "octopusdeploy_azure_cloud_service_deployment_targets" "data_lookup" {
  partial_name = "Azure"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_azure_cloud_service_deployment_targets.azure_cloud_service_deployment_targets[0].id
}