data "octopusdeploy_lifecycles" "example" {
  partial_name = "Simple"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_lifecycles.example.lifecycles[0].id
}