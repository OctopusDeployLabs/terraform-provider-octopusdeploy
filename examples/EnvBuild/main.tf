provider "octopusdeploy" {
  address  = var.serverURL
  apikey   = var.apiKey
  space_id = var.space
}

resource "octopusdeploy_environment" "newEnvironment" {
  description = "test environment"
  name        = var.environmentName
}

resource "octopusdeploy_lifecycle" "newLifecycle" {
  description = "test description"
  name        = var.lifecycleName

  release_retention_policy {
    quantity_to_keep = 30
    unit             = "Days"
  }

  depends_on = [octopusdeploy_environment.newEnvironment]
}

resource "octopusdeploy_project_group" "DevOpsProject" {
  description = "my test project group"
  name        = var.projectGroupName
}

resource "octopusdeploy_project" "DevOpsProject" {
  name             = var.projectName
  description      = "test project"
  lifecycle_id     = octopusdeploy_lifecycle.newLifecycle.id
  project_group_id = octopusdeploy_project_group.DevOpsProject.id

  depends_on = [octopusdeploy_project_group.DevOpsProject]
}

resource "octopusdeploy_aws_account" "aw" {
  name       = var.awsAccountName
  access_key = var.accessKey
  secret_key = var.secretKey
}

resource "octopusdeploy_azure_service_principal" "Azure" {
  client_id           = var.client_id
  key                 = var.key
  name                = "terratesttest"
  subscription_number = var.subscription_number
  tenant_id           = var.tenant_id
}

resource "octopusdeploy_channel" "newChannel" {
  description  = "test channel"
  lifecycle_id = octopusdeploy_lifecycle.newLifecycle.id
  name         = var.channelName
  project_id   = octopusdeploy_project.DevOpsProject.id
}

resource "octopusdeploy_variable" "newVariable" {
  name       = var.varName
  project_id = octopusdeploy_project.DevOpsProject.id
  type       = "String"
  value      = "testing123!@#"
}

output "octopus_deploy_environment" {
  value = octopusdeploy_environment.newEnvironment
}

output "octopus_deploy_lifecycle" {
  value = octopusdeploy_lifecycle.newLifecycle
}
