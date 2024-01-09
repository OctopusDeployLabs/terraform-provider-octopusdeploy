resource "octopusdeploy_offline_package_drop_deployment_target" "example" {
  environments                      = ["Environments-123", "Environment-321"]
  is_disabled                       = true
  machine_policy_id                 = "MachinePolicies-123"
  name                              = "Offline Package Drop Deployment Target (OK to Delete)"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
  thumbprint                        = "<thumbprint>"
  working_directory                 = "<working directory>"
  applications_directory            = "<applications directory>"
}
