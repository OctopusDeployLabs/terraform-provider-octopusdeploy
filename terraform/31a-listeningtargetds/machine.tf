data "octopusdeploy_listening_tentacle_deployment_targets" "listening_tentacle_deployment_targets" {
  partial_name    = "Test"
  skip            = 0
  take            = 1
}

output "data_lookup" {
  value = data.octopusdeploy_listening_tentacle_deployment_targets.listening_tentacle_deployment_targets.listening_tentacle_deployment_targets[0].id
}