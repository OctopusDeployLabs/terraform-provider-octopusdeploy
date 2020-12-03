data "octopusdeploy_project_groups" "example" {
  ids          = ["ProjectGroups-123", "ProjectGroups-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
