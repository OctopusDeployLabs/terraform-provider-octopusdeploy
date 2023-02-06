data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_azure_cloud_service_deployment_target" "target_azure" {
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  name                              = "Azure"
  roles                             = ["cloud"]
  account_id                        = "${octopusdeploy_azure_subscription_account.account_subscription.id}"
  cloud_service_name                = "servicename"
  storage_account_name              = "accountname"
  default_worker_pool_id            = ""
  health_status                     = "Unhealthy"
  is_disabled                       = false
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  shell_name                        = "Unknown"
  shell_version                     = "Unknown"
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []
  thumbprint                        = ""
  use_current_instance_count        = true
}