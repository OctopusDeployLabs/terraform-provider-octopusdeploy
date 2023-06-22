data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource "octopusdeploy_ssh_key_account" "account_ec2_sydney" {
  name                              = "ec2 sydney"
  description                       = ""
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  private_key_file                  = "${var.account_ec2_sydney_cert}"
  username                          = "ec2-user"
  private_key_passphrase            = "${var.account_ec2_sydney}"
}
variable "account_ec2_sydney" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The password associated with the certificate for account ec2 sydney"
}
variable "account_ec2_sydney_cert" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The certificate file for account ec2 sydney"
}

resource "octopusdeploy_ssh_connection_deployment_target" "target_3_25_215_87" {
  account_id            = "${octopusdeploy_ssh_key_account.account_ec2_sydney.id}"
  environments          = ["${octopusdeploy_environment.development_environment.id}"]
  fingerprint           = "d5:6b:a3:78:fa:fe:f5:ad:d4:79:4a:57:35:6a:32:ef"
  host                  = "3.25.215.87"
  name                  = "Test"
  roles                 = ["vm"]
  dot_net_core_platform = "linux-x64"
  machine_policy_id     = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
}