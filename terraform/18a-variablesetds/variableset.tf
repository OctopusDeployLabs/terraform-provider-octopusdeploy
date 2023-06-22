data "octopusdeploy_library_variable_sets" "lookup" {
  partial_name = "Test"
  skip = 0
  take = 1
}

output "data_lookup" {
  value = data.octopusdeploy_library_variable_sets.lookup.library_variable_sets[0].id
}