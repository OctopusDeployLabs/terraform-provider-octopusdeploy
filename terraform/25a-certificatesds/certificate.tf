data "octopusdeploy_certificates" "data_lookup" {
  archived     = false
  partial_name = "Test"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_certificates.data_lookup.certificates[0].id
}