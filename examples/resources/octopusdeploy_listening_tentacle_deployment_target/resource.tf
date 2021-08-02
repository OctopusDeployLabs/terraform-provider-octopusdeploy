resource "octopusdeploy_listening_tentacle_deployment_target" "example" {
  environments                      = ["Environments-123", "Environment-321"]
  is_disabled                       = true
  machine_policy_id                 = "MachinePolicy-123"
  name                              = "Listening Tentacle Deployment Target (OK to Delete)"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"
  tentacle_url                      = "https://example.com:1234/"
  thumbprint                        = "<thumbprint>"
}
