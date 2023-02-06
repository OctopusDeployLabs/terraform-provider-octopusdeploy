data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_listening_tentacle_deployment_target" "target_vm_listening_ngrok" {
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  name                              = "Test"
  roles                             = ["vm"]
  tentacle_url                      = "https://tentacle/"
  thumbprint                        = "55E05FD1B0F76E60F6DA103988056CE695685FD1"
  is_disabled                       = false
  is_in_process                     = false
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  shell_name                        = "Unknown"
  shell_version                     = "Unknown"
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []

  tentacle_version_details {
  }
}