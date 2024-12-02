resource "octopusdeploy_azure_openid_connect" "example" {
  name            = "Generic OpenID Connect Account (OK to Delete)"
  execution_subject_keys = ["space", "project"]
  audience = "api://default"
}