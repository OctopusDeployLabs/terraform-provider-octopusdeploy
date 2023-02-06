resource "octopusdeploy_token_account" "account_autopilot_service_account" {
  description                       = "A test account"
  name                              = "Token"
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  token                             = "secretgoeshere"
}