data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_offline_package_drop_deployment_target" "target_offline" {
  applications_directory            = "c:\\temp"
  working_directory                 = "c:\\temp"
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  name                              = "Test"
  roles                             = ["offline"]
  health_status                     = "Healthy"
  is_disabled                       = false
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  shell_name                        = "Unknown"
  shell_version                     = "Unknown"
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []
}