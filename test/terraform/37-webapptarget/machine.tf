data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_azure_web_app_deployment_target" "target_web_app" {
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  name                              = "Web App"
  roles                             = ["cloud"]
  account_id                        = "${octopusdeploy_azure_service_principal.account_sales_account.id}"
  resource_group_name               = "mattc-webapp"
  web_app_name                      = "mattc-webapp"
  health_status                     = "Unhealthy"
  is_disabled                       = false
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  shell_name                        = "Unknown"
  shell_version                     = "Unknown"
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []
  thumbprint                        = ""
  web_app_slot_name                 = "slot1"
}
resource "octopusdeploy_azure_service_principal" "account_sales_account" {
  name                              = "Sales Account"
  description                       = ""
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  application_id                    = "08a4a027-6f2a-4793-a0e5-e59a3c79189f"
  password                          = "${var.account_sales_account}"
  subscription_id                   = "3b50dcf4-f74d-442e-93cb-301b13e1e2d5"
  tenant_id                         = "3d13e379-e666-469e-ac38-ec6fd61c1166"
}
variable "account_sales_account" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The Azure secret associated with the account Sales Account"
}