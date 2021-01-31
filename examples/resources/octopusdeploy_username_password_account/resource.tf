resource "octopusdeploy_username_password_account" "example" {
  name     = "Username-Password Account (OK to Delete)"
  password = "###########" # get from secure environment/store
  username = "[username]"
}
