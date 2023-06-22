data "octopusdeploy_projects" "example" {
  partial_name           = "Test"
  skip                   = 0
  take                   = 1
}

output "data_lookup" {
  value = data.octopusdeploy_projects.example.projects[0].id
}