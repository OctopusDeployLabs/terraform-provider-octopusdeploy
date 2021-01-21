data "octopusdeploy_user_roles" "example" {
  ids          = ["UserRoles-123", "UserRoles-321"]
  partial_name = "Administra"
  skip         = 5
  take         = 100
}