---
page_title: "Creating a Resource"
subcategory: "Guides"
---

# Creating your first Octopus Deploy resource

## Provider Setup

The provider needs to be configured with proper credentials before it can be used, so first of all we need to create an API key in Octopus that we can then use in our Terraform config.

### Create an API key

* Log into the Octopus Web Portal, click your profile image and select **Profile**.
* Click **My API Keys**.
* Click **New API key**, state the purpose of the API key and click **Generate new**.
* Copy the new API key to your Terraform config file.

### Configure the provider in your Terraform config file.

```hcl
provider "octopusdeploy" {
  address = "https://octopus-cloud-or-server-uri"
  api_key = "octous_deploy_api_key"
}
```

## Creating your first Octopus Deploy resource

Once the provider is configured, you can apply the Octopus Deploy resources defined in your Terraform config file.

The following is an example Terraform config file that creates a new `Project Group` in the `Default` space.

```hcl
terraform {
    required_providers {
        octopusdeploy = {
            source  = OctopusDeployLabs/octopusdeploy
            version = ">= 0.7.64"
        }
    }
}

provider "octopusdeploy" {
  address = "https://octopus-cloud-or-server-uri"
  api_key = "octous_deploy_api_key"
}

resource "octopusdeploy_project_group" "DevOpsProjects" {
  name        = "DevOps Projects"
  description = "My DevOps projects group"
}
```

1. Use `terraform init` to download the specified version of the Octopus Deploy provider.
2. Next, use `terraform plan` to display a list of resources to be created, and highlight any possible unknown attributes at apply time.
3. Finally, use `terraform apply` to create the resource shown in the output from the step above.