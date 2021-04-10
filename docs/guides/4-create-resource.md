---
page_title: "4. Creating a Resource"
subcategory: "Guides"
---

# 4. Creating a Resource

Now that the environment is up and running, it's time to create a resource to test. Just so we can see how it works, we'll create one resource.

Before we start creating a Terraform configuration, we'll need a few components.

1. The server URL, which is going to be `localhost:8080`.
2. An API key which is generated from the Octopus Deploy server.

When connecting to the Octopus Deploy Terraform provider, the server URL and API key is needed for authentication. That way, Terraform knows what environment to connect to.

Once you have the API key, you can move on to the next step.

## main.tf

This file is where you put the newly-created resource. For example, you can create a resource to create a new project group:

```terraform
provider "octopusdeploy" {
  address = "localhost:8080"
  apikey  = api_key
  space   = "Default"
}

resource "octopusdeploy_project_group" "DevOpsProject" {
  name        = "testProject"
  description = "my test project group"
}
```

With the above HCL code, the Terraform resource can be created with the standard Terraform commands:

1. `terraform init`
2. `terraform plan`
3. `terraform apply`