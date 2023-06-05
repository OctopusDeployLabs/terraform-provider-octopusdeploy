resource "octopusdeploy_aws_account" "account_aws_account" {
  name                              = "AWS Account"
  description                       = ""
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  access_key                        = "ABCDEFGHIJKLMNOPQRST"
  secret_key                        = "secretgoeshere"
}
