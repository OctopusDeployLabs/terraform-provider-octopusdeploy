data "octopusdeploy_teams" "teams" {
  partial_name = "Deployers"
  skip = 0
  take = 1
}

output "teams_lookup" {
  value = data.octopusdeploy_teams.teams.teams[0].id
}


data "octopusdeploy_user_roles" "example" {
  partial_name = "Project Deployer"
  skip         = 0
  take         = 1
}

output "roles_lookup" {
  value = data.octopusdeploy_user_roles.example.user_roles[0].id
}

data "octopusdeploy_users" "example" {
  filter  ="Bob Smith"
  skip = 0
  take = 1
}

output "users_lookup" {
  value = data.octopusdeploy_users.example.users[0].id
}