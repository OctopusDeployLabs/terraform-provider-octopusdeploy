resource "octopusdeploy_azure_openid_connect" "example" {
  name            = "Generic OpenID Connect Account (OK to Delete)"
  execution_subject_keys = ["space", "project"]
  health_subject_keys = ["space", "target", "type"]
  account_test_subject_keys = ["space", "type"]
  audience = "api://Default"
}