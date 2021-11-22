data "octopusdeploy_script_modules" "example" {
  ids          = ["LibraryVariableSets-123", "LibraryVariableSets-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
