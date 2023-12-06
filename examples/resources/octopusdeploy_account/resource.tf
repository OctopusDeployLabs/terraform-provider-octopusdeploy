# create an Amazon web services account
resource "octopusdeploy_account" "amazon_web_services_account" {
  access_key   = "access-key"
  account_type = "AmazonWebServicesAccount"
  name         = "AWS Account (OK to Delete)"
  secret_key   = "###########" # required; get from secure environment/store
}

resource "octopusdeploy_account" "amazon_web_services_openid_connect_account" {
  account_type = "AwsOIDCAccount"
  name         = "AWS OIDC Account (OK to Delete)"
  role_arn = "arn:aws:iam::sourceAccountId:roleroleName"
  session_duration = "3600"
  execution_subject_keys = ["space", "project"]
  health_subject_keys = ["space", "target"]
  account_test_subject_keys = ["space", "type"]
}

# create an Azure service principal account
resource "octopusdeploy_account" "azure_service_principal_account" {
  account_type    = "AzureServicePrincipal"
  application_id  = "00000000-0000-0000-0000-000000000000"
  name            = "Azure Service Principal Account (OK to Delete)"
  password        = "###########" # required; get from secure environment/store
  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
}

# create an Azure OpenID Connect account
resource "octopusdeploy_account" "azure_openid_connect_account" {
  account_type    = "AzureOIDC"
  application_id  = "00000000-0000-0000-0000-000000000000"
  name            = "Azure OpenID Connect Account (OK to Delete)"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
  execution_subject_keys = ["space", "project"]
  health_subject_keys = ["space", "target"]
  account_test_subject_keys = ["space", "type"]
}

# create an Azure subscription account
resource "octopusdeploy_account" "azure_subscription_account" {
  account_type    = "AzureSubscription"
  name            = "Azure Subscription Account (OK to Delete)"
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

# create a SSH key pair account
resource "octopusdeploy_account" "ssh_key_pair_account" {
  account_type     = "SshKeyPair"
  name             = "SSH Key Pair Account (OK to Delete)"
  private_key_file = "[private_key_file]"
  username         = "[username]"
}

# create a username-password account
resource "octopusdeploy_account" "username_password_account" {
  account_type = "UsernamePassword"
  name         = "Username-Password Account (OK to Delete)"
  username     = "[username]"
}

# create a token account
resource "octopusdeploy_account" "token_account" {
  account_type = "Token"
  name         = "Token Account (OK to Delete)"
  token        = "[token]"
}
