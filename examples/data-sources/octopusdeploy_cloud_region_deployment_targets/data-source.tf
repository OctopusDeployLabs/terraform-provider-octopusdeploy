data "octopusdeploy_cloud_region_deployment_targets" "example" {
  deployment_id = "Defau"
  environments  = ["Environments-123", "Environments-321"]
  ids           = ["Machines-123"]
  skip          = 5
  take          = 100
}
