---
page_title: "5. Creating Multiple Resources"
subcategory: "Guides"
---

# 5. Creating Multiple Resources

In step 4, you learned about creating one resource. That was just to test the waters so you could see how the provider would interact with Octopus Deploy. Now that you've done so, it's time to really test the waters, kick it into 6th gear, and take this puppy for a spin.

First, you'll set up a `main.tf` configuration with multiple resources to create.

## Main.tf
Below is the `main.tf` configuration to use.

The below configuration will create:

1. An Octopus Deploy environment
2. Lifecycle
3. ProjectGroup
4. Project 
5. AWS Account
6. Azure Account
7. Channel

To test the output and see some resources that were created, the terminal will show at the end:

1. The output of the new environment
2. The output of the new lifecycle

```hcl
provider "octopusdeploy" {
  address = var.serverURL
  apikey  = var.apiKey
  space_id = var.space
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

output "octopus_deploy_environment" {
  value = octopusdeploy_environment.newEnvironment
}

output "octopus_deploy_lifecycle" {
  value = octopusdeploy_lifecycle.newLifecycle
}

```

## Variables.tf
Because you don't want to have to add in static values manually into the larger Terraform configuration, you'll set up some variables.

All variables will have a type of `String` and won't have default values because you'll be passing in the values at runtime using a `terraform.tfvars` configuration file.

```hcl
variable "serverURL" {
    type = string
}

variable "apiKey" {
    type = string
}

variable "space" {
    type = string
}

variable "accessKey" {
    type = string
}

variable "secretKey" {
    type = string
}

variable "channelName" {
    type = string
}

variable "environmentName" {
    type = string
}

variable "lifecycleName" {
    type = string
}

variable "projectGroupName" {
    type = string
}

variable "projectName" {
    type = string
}

variable "client_id" {
    type = string
}

variable "tenant_id" {
    type = string
}

variable "subscription_number" {
    type = string
}

variable "key" {
    type = string
}

variable "awsAccountName" {
    type = string
}
```

## Terraform.tfvars
Finally, to pass in new values at runtime. You'll have a `terraform.tfvars` configuration file. 

The values that you see below are placeholders. Feel free to add in your own values. If you don't have values certain parts, for example, the Azure Service Principle which requires a GUID, simply add in a random GUID.

```hcl
serverURL = "https://octopus-cloud-or-server-uri"
apiKey = "octous_deploy_api_key"
space = "Default"
accessKey = "aws_access_key"
secretKey = "aws_secret_key"
channelName = "DevChannel"
environmentName = "DevEnv"
lifecycleName = "DevLifecycle"
projectGroupName = "DevTest"
projectName = "DevTest1"
client_id = "azure_app_registration_client_id"
tenant_id = "azure_app_registration_tenant_id"
subscription_number = "azure_subscription_id"
key = "azure_app_registration_access_key"
awsAccountName = "DevAWS"
```

When you're all done, the directory containing the configurations should look like the screenshot below.

![Build](images/build.png)