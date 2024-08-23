data "octopusdeploy_kubernetes_agent_workers" "all_workers" {
}

output "data_lookup_kubernetes_worker_1_id" {
  value = data.octopusdeploy_kubernetes_agent_workers.all_workers.kubernetes_agent_workers[0].id
}

output "data_lookup_kubernetes_worker_2_id" {
  value = data.octopusdeploy_kubernetes_agent_workers.all_workers.kubernetes_agent_workers[1].id
}