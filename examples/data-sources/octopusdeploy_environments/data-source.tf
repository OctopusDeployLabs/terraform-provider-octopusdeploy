data "octopusdeploy_environments" "example" {
  ids          = ["Environments-123", "Environments-321"]
  name         = "Production"
  partial_name = "Produc"
  skip         = 5
  take         = 100
}
