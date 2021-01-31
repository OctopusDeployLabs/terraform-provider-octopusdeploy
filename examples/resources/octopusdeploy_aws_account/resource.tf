resource "octopusdeploy_aws_account" "example" {
  access_key   = "access-key"
  name         = "AWS Account (OK to Delete)"
  secret_key   = "###########" # required; get from secure environment/store
}