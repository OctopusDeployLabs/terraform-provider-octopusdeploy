# create an Amazon web services account variable
resource "octopusdeploy_variable" "amazon_web_services_account_variable" {
  owner_id  = "Projects-123"
  type      = "AmazonWebServicesAccount"
  name      = "My AWS Account (OK to Delete)"
  value     = "Accounts-123"
}

# create an Azure service principal account variable
resource "octopusdeploy_account" "azure_service_principal_account_variable" {
  owner_id  = "Projects-123"
  type      = "AzureAccount"
  name      = "My Azure Service Principal Account (OK to Delete)"
  value     = "Accounts-123"
}

# create a Google Cloud account variable
resource "octopusdeploy_variable" "google_cloud_account_variable" {
  owner_id  = "Projects-123"
  type      = "GoogleCloudAccount"
  name      = "My Google Cloud Account (OK to Delete)"
  value     = "Accounts-123"
}

# create a Certificate variable
resource "octopusdeploy_variable" "certificate_variable" {
  owner_id  = "Projects-123"
  type      = "Certificate"
  name      = "My Certificate (OK to Delete)"
  value     = "Certificates-123"
}

# create a Sensitive variable
resource "octopusdeploy_variable" "sensitive_variable" {
  owner_id        = "Projects-123"
  type            = "Sensitive"
  name            = "My Sensitive Value (OK to Delete)"
  is_sensitive    = true
  sensitive_value = "YourSecrets"
}

# create a String variable
resource "octopusdeploy_variable" "string_variable" {
  owner_id  = "Projects-123"
  type      = "String"
  name      = "My String Value (OK to Delete)"
  value     = "PlainText"
}

# create a WorkerPool variable
resource "octopusdeploy_variable" "workerpool_variable" {
  owner_id  = "Projects-123"
  type      = "WorkerPool"
  name      = "My WorkerPool (OK to Delete)"
  value     = "WorkerPools-123"
}

# create a prompted variable
resource "octopusdeploy_variable" "prompted_variable" {
  owner_id  = "Projects-123"
  type      = "String"
  name      = "My Prompted Variable (OK to Delete)"
  prompt {
    description = "Variable Description"
    is_required = true
    label       = "Variable Label"
  }
}
