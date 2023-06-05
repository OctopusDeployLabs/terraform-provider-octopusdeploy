data "octopusdeploy_cloud_region_deployment_targets" "example" {
  partial_name = "Test"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_cloud_region_deployment_targets.example.cloud_region_deployment_targets[0].id
}