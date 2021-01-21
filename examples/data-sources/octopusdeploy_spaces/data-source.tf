data "octopusdeploy_spaces" "spaces" {
  ids          = ["Spaces-123", "Spaces-321"]
  name         = "Default"
  partial_name = "Defau"
  skip         = 5
  take         = 100
}