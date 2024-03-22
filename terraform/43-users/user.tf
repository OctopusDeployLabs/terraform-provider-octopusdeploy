resource "octopusdeploy_user" "serviceaccount" {
  display_name  = "Service Account"
  email_address = "a@a.com"
  is_active     = true
  is_service    = true
  username      = "saccount"
}

resource "octopusdeploy_user" "deployer" {
  display_name  = "Bob Smith"
  email_address = "bob.smith@example.com"
  is_active     = false
  is_service    = false
  username      = "bsmith"
  password      = "Password01!"
}

resource "octopusdeploy_team" "deployers" {
  name  = "Deployers"
  users = [octopusdeploy_user.deployer.id]
}

resource "octopusdeploy_scoped_user_role" "deploy" {
  space_id     = var.octopus_space_id
  team_id      = octopusdeploy_team.deployers.id
  user_role_id = "userroles-projectdeployer"
}

resource "octopusdeploy_scoped_user_role" "release" {
  space_id     = var.octopus_space_id
  team_id      = octopusdeploy_team.deployers.id
  user_role_id = "userroles-releasecreator"
}