resource "octopusdeploy_ssh_connection_deployment_target" "example" {
  name        = "SSH Connection Deployment Target (OK to Delete)"
  fingerprint = "[fingerprint]"
  host        = "[host]"
  port        = 22
}
