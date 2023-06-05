data "octopusdeploy_script_modules" "example" {
  partial_name = "Test2"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_script_modules.example.script_modules[0].id
}