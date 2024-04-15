resource "octopusdeploy_username_password_account" "account_user_pass" {
  description                       = "A test account"
  name                              = "Username Password"
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  username                          = "admin"
  password                          = "secretgoeshere"
}