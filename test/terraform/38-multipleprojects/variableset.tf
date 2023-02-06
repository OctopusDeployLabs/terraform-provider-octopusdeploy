resource "octopusdeploy_library_variable_set" "octopus_library_variable_set" {
  name = "Test"
  description = "Test variable set"
}

resource "octopusdeploy_library_variable_set" "octopus_library_variable_set2" {
  name = "Test2"
  description = "Test variable set"
}

resource "octopusdeploy_variable" "octopus_admin_api_key" {
  name = "Test.Variable"
  type = "String"
  description = "Test variable"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = "test"
}