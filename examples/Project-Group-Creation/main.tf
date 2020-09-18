provider "octopusdeploy" {
    address = var.serverURL
    apikey  = var.apiKey
    space_id   = var.space
}

resource "octopusdeploy_project_group" "DevOpsProject" {
    name = "testProject"
    description = "my test project group"
}
