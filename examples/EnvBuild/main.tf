provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space   = var.space
}

resource "octopusdeploy_environment" "newEnvironment" {
  name        = var.environmentName
  description = "test environment"
}

resource "octopusdeploy_lifecycle" "newLifecycle" {
  name        = var.lifecycleName
  description = "test description"
  release_retention_policy {
    quantity_to_keep = 30
    unit             = "Days"
  }

  depends_on = [octopusdeploy_environment.newEnvironment]
}

resource "octopusdeploy_project_group" "DevOpsProject" {
  name        = var.projectGroupName
  description = "my test project group"
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
  name                = "terratesttest"
  client_id           = var.client_id
  tenant_id           = var.tenant_id
  subscription_number = var.subscription_number
  key                 = var.key
}

resource "octopusdeploy_channel" "newChannel" {
  name         = var.channelName
  description  = "test channel"
  lifecycle_id = octopusdeploy_lifecycle.newLifecycle.id
  project_id   = octopusdeploy_project.DevOpsProject.id
}

resource "octopusdeploy_variable" "newVariable" {
  name       = var.varName
  value      = "testing123!@#"
  project_id = octopusdeploy_project.DevOpsProject.id
  type       = "String"
}

output "octopus_deploy_environment" {
  value = octopusdeploy_environment.newEnvironment
}

output "octopus_deploy_lifecycle" {
  value = octopusdeploy_lifecycle.newLifecycle
}
