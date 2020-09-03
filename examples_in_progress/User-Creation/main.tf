provider "octopusdeploy" {
    address = var.serverURL
    apikey  = var.apiKey
    space   = var.space
}

resource "octopusdeploy_user" "NewUser" {
    UserName = var.userName
    DisplayName = var.displayName
}
