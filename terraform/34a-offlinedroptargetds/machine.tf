data "octopusdeploy_offline_package_drop_deployment_targets" "lookup" {
  partial_name = "Test"
  skip = 0
  take = 1
}

output "data_lookup" {
  value = data.octopusdeploy_offline_package_drop_deployment_targets.lookup.offline_package_drop_deployment_targets[0].id
}