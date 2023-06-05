data "octopusdeploy_tenants" "tenants" {
  partial_name = "Team A"
  skip = 0
  take = 1
}

output "tenants_lookup" {
  value = data.octopusdeploy_tenants.tenants.tenants[0].id
}