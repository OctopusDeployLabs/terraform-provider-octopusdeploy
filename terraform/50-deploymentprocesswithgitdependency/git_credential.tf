resource "octopusdeploy_git_credential" "test_git_credential" {
  description = "test git credential"
  name        = "test"
  type        = "UsernamePassword"
  username    = "admin"
  password    = "Password01!"
}