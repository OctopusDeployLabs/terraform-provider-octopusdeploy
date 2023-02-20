resource "octopusdeploy_gcp_account" "account_google" {
  description                       = "A test account"
  name                              = "Google"
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  json_key                          = "secretgoeshere"
}