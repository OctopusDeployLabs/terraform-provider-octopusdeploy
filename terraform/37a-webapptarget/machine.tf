data "octopusdeploy_azure_web_app_deployment_targets" "data_lookup" {
  partial_name = "Web App"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_azure_web_app_deployment_targets.data_lookup.azure_web_app_deployment_targets[0].id
}