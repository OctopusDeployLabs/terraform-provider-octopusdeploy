resource "octopusdeploy_ssh_connection_deployment_target" "example" {
  name                              = "SSH Connection Deployment Target (OK to Delete)"
  fingerprint                       = "<fingerprint>"
  host                              = "<host>"
  port                              = 22
  environments                      = ["Environments-123", "Environment-321"]
  account_id                        = "Accounts-1"
  roles                             = ["Development Team", "System Administrators"]
}
