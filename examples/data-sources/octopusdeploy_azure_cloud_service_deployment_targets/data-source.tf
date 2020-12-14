data "octopusdeploy_azure_cloud_service_deployment_targets" "example" {
  health_statuses = ["Healthy", "Unavailable"]
  ids             = ["Machines-123", "Machines-321"]
  partial_name    = "Defau"
  skip            = 5
  take            = 100
}
