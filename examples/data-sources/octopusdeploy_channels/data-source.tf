data "octopusdeploy_channels" "example" {
  ids          = ["Channels-123", "Channels-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
