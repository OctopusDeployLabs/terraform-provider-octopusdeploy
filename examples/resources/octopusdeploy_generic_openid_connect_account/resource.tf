resource "octopusdeploy_generic_openid_connect_account" "example" {
  name            = "Generic OpenID Connect Account (OK to Delete)"
  execution_subject_keys = ["space", "project"]
  audience = "api://default"
}