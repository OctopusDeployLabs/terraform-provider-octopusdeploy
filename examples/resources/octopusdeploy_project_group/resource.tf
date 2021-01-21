resource "octopusdeploy_project_group" "example" {
  description  = "The development project group."
  environments = ["Environments-123", "Environments-321"]
  name         = "Development Project Group (OK to Delete)"
}
