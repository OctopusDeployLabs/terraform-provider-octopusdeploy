resource "octopusdeploy_account" "amazon_web_services_account" {
  access_key   = "access-key"
  account_type = "AmazonWebServicesAccount"
  name         = "AWS Account (OK to Delete)"
  secret_key   = "###########" # required; get from secure environment/store
}