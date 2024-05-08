resource "octopusdeploy_polling_subscription_id" "base" {}
resource "octopusdeploy_polling_subscription_id" "optionals" {
  dependencies = {
    "base_id" = octopusdeploy_polling_subscription_id.base.id
  }
}

output "base_id" {
  value = octopusdeploy_polling_subscription_id.base.id
}

output "base_polling_uri" {
  value = octopusdeploy_polling_subscription_id.base.polling_uri
}