data "octopusdeploy_channels" "data_lookup" {
  partial_name = "Test"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_channels.data_lookup.channels[0].id
}