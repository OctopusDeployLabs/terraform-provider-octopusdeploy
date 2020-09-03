provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_usernamepassword_account" "username" {
    name = "newuser"
    username = "testing"
    password = "testing123!@#"
}