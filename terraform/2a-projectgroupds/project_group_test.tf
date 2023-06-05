data "octopusdeploy_project_groups" "example" {
  partial_name = "Test"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_project_groups.example.project_groups[0].id
}