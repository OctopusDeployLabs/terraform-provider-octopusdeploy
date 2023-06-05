provider "octopusdeploy" {
  address = "${var.octopus_server}"
  api_key = "${var.octopus_apikey}"
}
