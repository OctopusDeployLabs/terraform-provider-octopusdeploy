data "octopusdeploy_cloud_region_deployment_targets" "example" {
  environments = ["Environments-123", "Environments-321"]
  ids          = ["Machines-123"]
  name         = "Azure North America"
  partial_name = "Azure Nor"
  skip         = 5
  take         = 100
}
