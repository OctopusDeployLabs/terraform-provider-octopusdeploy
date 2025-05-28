---
page_title: "Creating Multiple Resources"
subcategory: "Guides"
---

# Creating multiple Octopus Deploy resources

In Creating a Resource guide, you learned about creating one resource. In this guide we will create multiple resources with the same Terraform config file.

The below Terraform configuration will create the following Octopus Deploy resources:
* Environment
* Lifecycle
* Project Group
* Project 
* AWS Account
* Azure Account
* Channel

To test the output and see some resources that were created, the terminal will show at the end:
* The output of the new environment
* The output of the new lifecycle

```hcl
terraform {
    required_providers {
        octopusdeploy = {
            source = "octopus.com/com/octopusdeploy"
            version = ">= 0.7.64"
        }
    }
}

provider "octopusdeploy" {
  address  = var.serverURL
  api_key  = var.apiKey
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

  depends_on  = [octopusdeploy_environment.newEnvironment]
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

  depends_on  = [octopusdeploy_project_group.DevOpsProject]
}

resource "octopusdeploy_aws_account" "aw" {
  name       = var.awsAccountName
  access_key = var.accessKey
  secret_key = var.secretKey
}

resource "octopusdeploy_azure_service_principal" "Azure" {
  name            = "terratesttest"
  application_id  = var.application_id
  tenant_id       = var.tenant_id
  subscription_id = var.subscription_id
  password        = var.password
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

Because you don't want to have to add in static values manually into the larger Terraform configuration, you'll set up some variables in a file named `variables.tf`.

_All variables will have a type of `String` and won't have default values because you'll be passing in the values at runtime using a `terraform.tfvars` configuration file._

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

variable "application_id" {
    type = string
}

variable "tenant_id" {
    type = string
}

variable "subscription_id" {
    type = string
}

variable "password" {
    type = string
}

variable "awsAccountName" {
    type = string
}
```

Finally, to pass in new values at runtime. You'll create a `terraform.tfvars` configuration file. 

_The values that you see below are placeholders. Feel free to add in your own values. If you don't have values certain parts, for example, the Azure Service Principle which requires a GUID, simply add in a random GUID._

```hcl
serverURL        = "https://octopus-cloud-or-server-uri"
apiKey           = "octous_deploy_api_key"
space            = "Spaces-1"
accessKey        = "aws_access_key"
secretKey        = "aws_secret_key"
channelName      = "DevChannel"
environmentName  = "DevEnv"
lifecycleName    = "DevLifecycle"
projectGroupName = "DevTest"
projectName      = "DevTest1"
application_id   = "azure_app_registration_client_id"
tenant_id        = "azure_app_registration_tenant_id"
subscription_id  = "azure_subscription_id"
password         = "azure_app_registration_access_key"
awsAccountName   = "DevAWS"
```

When you're all done, your directory containing the configurations should look like the below folder structure.

```
terraform_example/
├── main.tf
├── .terraform
│   └── providers
│       └── octopus.com
│           └── com
│               └── octopusdeploy
│                   └── 0.7.64
├── .terraform.lock.hcl
├── terraform.tfstate
├── terraform.tfvars
└── variables.tf
```