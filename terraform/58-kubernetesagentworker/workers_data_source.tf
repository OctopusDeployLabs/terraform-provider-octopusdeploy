data "octopusdeploy_kubernetes_agent_workers" "all_workers" {
}

output "data_lookup" {
  value = data.octopusdeploy_kubernetes_agent_workers.all_workers
}