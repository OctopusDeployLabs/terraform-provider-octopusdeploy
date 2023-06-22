data "octopusdeploy_environments" "data_lookup" {
  partial_name = "Development"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_environments.data_lookup.environments[0].id
}