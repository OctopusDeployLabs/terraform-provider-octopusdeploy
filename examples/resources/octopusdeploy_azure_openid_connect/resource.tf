resource "octopusdeploy_azure_openid_connect" "example" {
  application_id  = "00000000-0000-0000-0000-000000000000"
  name            = "Azure OpenID Connect Account (OK to Delete)"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  execution_subject_keys = ["space", "project"]
  health_subject_keys = ["space", "target", "type"]
  account_test_subject_keys = ["space", "type"]
  audience = "api://AzureADTokenExchange"
}