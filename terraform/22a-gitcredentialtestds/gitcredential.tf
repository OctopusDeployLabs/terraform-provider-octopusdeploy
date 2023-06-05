data "octopusdeploy_git_credentials" "data_lookup" {
  name = "test"
  skip = 0
  take = 1
}

output "data_lookup" {
  value = data.octopusdeploy_git_credentials.data_lookup.git_credentials[0].id
}