resource "octopusdeploy_azure_service_principal" "account_azure" {
  description                       = "Azure Account"
  name                              = "Azure"
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  application_id                    = "2eb8bd13-661e-489c-beb9-4103efb9dbdd"
  password                          = "secretgoeshere"
  subscription_id                   = "95bf77d2-64b1-4ed2-9de1-b5451e3881f5"
  tenant_id                         = "18eb006b-c3c8-4a72-93cd-fe4b293f82ee"
}