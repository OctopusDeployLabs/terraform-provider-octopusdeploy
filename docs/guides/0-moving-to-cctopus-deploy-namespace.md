---
page_title: "0. Moving to Octopus Deploy Namespace"
subcategory: "Guides"
---

# 0. Moving to Octopus Deploy Namespace

The aim of this guide is to help move your pre-existing *OctopusDeployLabs* provider configuration to the *OctopusDeploy* namespace while maintaining state.

> Keep in mind, it's important to finish all steps within guide before running `terraform apply`.

## 1. Update the Provider Block

Change from the *OctopusDeployLabs* provider source to the *OctopusDeploy* provider source.

Before:

```terraform
terraform {
  required_providers {
    octopusdeploy = {
      source = "OctopusDeployLabs/octopusdeploy"
      version = "0.43.x"
    }
  }
}
```

After:

```terraform
terraform {
  required_providers {
    octopusdeploy = {
      source = "OctopusDeploy/octopusdeploy"
      version = "1.x.x"
    }
  }
}
```

If copying from the example above please ensure to set the latest version.

## 2. Pull the Provider

Run the following to tell terraform to pull the provider under the *OctopusDeploy* namespace.

```shell
terraform init -upgrade
```

## 3. Move the Existing State

This is the key step to maintain state.
Tell terraform to map the *OctopusDeployLabs* namespace state to the *OctopusDeploy* namespace.

```shell
terraform state replace-provider OctopusDeployLabs/octopusdeploy OctopusDeploy/octopusdeploy
```

## 4. Verify

To verify the resources have moved to the *OctopusDeploy* namespace correctly run the following.

```shell
terraform plan
```

The plan should show no unexpected changes.
If the version being upgrading to introduced some breaking changes you may see some changes as expected.

## 5. Done

Success! Continue to use the provider as normal.
