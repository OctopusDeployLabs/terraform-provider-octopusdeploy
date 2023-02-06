data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_polling_tentacle_deployment_target" "target_desktop_3e4k4r8" {
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  name                              = "Test"
  roles                             = ["vm"]
  tentacle_url                      = "poll://abcdefghijklmnopqrst/"
  is_disabled                       = false
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  shell_name                        = "PowerShell"
  shell_version                     = "5.1.22621.1"
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []

  tentacle_version_details {
  }

  thumbprint = "1854A302E5D9EAC1CAA3DA1F5249F82C28BB2B86"
}