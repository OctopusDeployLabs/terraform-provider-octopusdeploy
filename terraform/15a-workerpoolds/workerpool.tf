data "octopusdeploy_worker_pools" "example" {
  partial_name = "Docker"
  skip = 0
  take = 1
}

output "data_lookup" {
  value = data.octopusdeploy_worker_pools.example.worker_pools[0].id
}