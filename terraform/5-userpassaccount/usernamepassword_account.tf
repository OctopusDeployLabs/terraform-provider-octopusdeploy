resource "octopusdeploy_username_password_account" "account_gke" {
  description                       = "A test account"
  name                              = "GKE"
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  username                          = "admin"
  password                          = "secretgoeshere"
}