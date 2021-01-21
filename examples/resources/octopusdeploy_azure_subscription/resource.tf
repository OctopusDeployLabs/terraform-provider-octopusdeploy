resource "octopusdeploy_account" "azure_subscription_account" {
  account_type    = "AzureSubscription"
  name            = "Azure Subscription Account (OK to Delete)"
  subscription_id = "00000000-0000-0000-0000-000000000000"
}