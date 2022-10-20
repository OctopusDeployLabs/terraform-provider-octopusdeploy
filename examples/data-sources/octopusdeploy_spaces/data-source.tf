data "octopusdeploy_spaces" "spaces" {
  ids          = ["Spaces-123", "Spaces-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
