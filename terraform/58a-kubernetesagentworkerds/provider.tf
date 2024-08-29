provider "octopusdeploy" {
  address  = "${var.octopus_server_58-kubernetesagentworker}"
  api_key  = "${var.octopus_apikey_58-kubernetesagentworker}"
  space_id = "${var.octopus_space_id}"
}
