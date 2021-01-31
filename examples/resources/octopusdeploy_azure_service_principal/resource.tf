resource "octopusdeploy_azure_service_principal" "example" {
  application_id  = "00000000-0000-0000-0000-000000000000"
  name            = "Azure Service Principal Account (OK to Delete)"
  password        = "###########" # required; get from secure environment/store
  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
}