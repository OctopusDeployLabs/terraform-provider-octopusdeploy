data "octopusdeploy_projects" "example" {
  cloned_from_project_id = "Projects-456"
  ids                    = ["Projects-123", "Projects-321"]
  is_clone               = true
  name                   = "Default"
  partial_name           = "Defau"
  skip                   = 5
  take                   = 100
}
