data "octopusdeploy_kubernetes_agent_deployment_targets" "test" {
  take         = 1
  skip         = 0
  partial_name = "minimum-agent"
}

output "data_lookup" {
  value = data.octopusdeploy_kubernetes_agent_deployment_targets.test.kubernetes_agent_deployment_targets[0].id
}