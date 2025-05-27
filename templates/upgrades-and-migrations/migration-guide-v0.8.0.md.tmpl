---
page_title: "Migrating to v0.8.0"
subcategory: "Upgrades & Migrations"
---

`v0.8.0` includes a number of breaking changes. This guide provides an outline of these changes along with tips on migration.

This guide pertains to configuration that references one of the following resources:

* `octopusdeploy_account` (was deprecated; now removed)
* `octopusdeploy_deployment_target` (was deprecated; now removed)
* `octopusdeploy_feed` (was deprecated; now removed)
* `octopusdeploy_tag_set` (changed)
* `octopusdeploy_tag` (new; added)

As always, please ensure to [`validate`](https://www.terraform.io/cli/commands/validate) your configuration and review the changes from [`plan`](https://www.terraform.io/cli/commands/plan) before committing changes through [`apply`](https://www.terraform.io/cli/commands/apply).

## Pinning a Provider Configuration to v0.7.73 (or Earlier)

If the configuration fails to [`validate`](https://www.terraform.io/cli/commands/validate) after updating to `v0.8.0` then you can pin the version of the provider to the previous version: `v0.7.73`:

```terraform
terraform {
  required_providers {
    octopusdeploy = {
      source = "OctopusDeployLabs/octopusdeploy"
      version = "0.7.73" # avoid qualifiers for versions, which can be dangerous until v1.0.0
    }
  }
}
```

At a later time, you may migrate your existing state and configuration to the new resource types.

## Tag Sets and Tags

In `v0.8.0`, the `octopusdeploy_tag_set` has been modified by hoisting its embedded `tag` blocks into a separate a distinct resource (`octopusdeploy_tag`).

Prior to `v0.8.0`, a tag set and tag(s) could be defined as follows:

```terraform
resource "octopusdeploy_tag_set" "test-tag-set" {
  name = "Test Tag Set (OK to Delete)"
 
  tag {
      color = "#FF0000"
      name  = "test-tag-1"
  }

  tag {
      color = "#00FF00"
      name  = "test-tag-2"
  }
}
```

In `v0.8.0`, the schema has been modified to resemble this:

```terraform
resource "octopusdeploy_tag_set" "test-tag-set" {
  name = "Test Tag Set (OK to Delete)"
}

resource "octopusdeploy_tag" "us-west" {
  color      = "#FF0000"
  name       = "test-tag-1"
  tag_set_id = octopusdeploy_tag_set.test-tag-set.id
}

resource "octopusdeploy_tag" "us-east" {
  color      = "#00FF00"
  name       = "test-tag-2"
  tag_set_id = octopusdeploy_tag_set.test-tag-set.id
}
```

The `octopusdeploy_tag` resource is new and has a required property, `tag_set_id` which associates it with an `octopusdeploy_tag_set` resource.

## octopusdeploy_account (removed)

The `octopusdeploy_account` was marked as deprecated and it has removed in `v0.8.0`. There are equivalent replacements available that provide more robust validation:

* [`octopusdeploy_aws_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/aws_account)
* [`octopusdeploy_azure_service_principal`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/azure_service_principal)
* [`octopusdeploy_azure_subscription_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/azure_subscription_account)
* [`octopusdeploy_gcp_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/gcp_account)
* [`octopusdeploy_ssh_key_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/ssh_key_account)
* [`octopusdeploy_token_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/token_account)
* [`octopusdeploy_username_password_account`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/username_password_account)

## octopusdeploy_deployment_target (removed)

The `octopusdeploy_deployment_target` was marked as deprecated and it has removed in `v0.8.0`. There are equivalent replacements available that provide more robust validation:

* [`octopusdeploy_azure_cloud_service_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/azure_cloud_service_deployment_target)
* [`octopusdeploy_azure_service_fabric_cluster_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/azure_service_fabric_cluster_deployment_target)
* [`octopusdeploy_azure_web_app_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/azure_web_app_deployment_target)
* [`octopusdeploy_cloud_region_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/cloud_region_deployment_target)
* [`octopusdeploy_kubernetes_cluster_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/kubernetes_cluster_deployment_target)
* [`octopusdeploy_listening_tentacle_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/listening_tentacle_deployment_target)
* [`octopusdeploy_offline_package_drop_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/offline_package_drop_deployment_target)
* [`octopusdeploy_polling_tentacle_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/polling_tentacle_deployment_target)
* [`octopusdeploy_ssh_connection_deployment_target`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/ssh_connection_deployment_target)

## octopusdeploy_feed (removed)

The `octopusdeploy_feed` was marked as deprecated and it has removed in `v0.8.0`. There are equivalent replacements available that provide more robust validation:

* [`octopusdeploy_aws_elastic_container_registry`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/aws_elastic_container_registry)
* [`octopusdeploy_docker_container_registry`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/docker_container_registry)
* [`octopusdeploy_github_repository_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/github_repository_feed)
* [`octopusdeploy_helm_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/helm_feed)
* [`octopusdeploy_maven_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/maven_feed)
* [`octopusdeploy_nuget_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/nuget_feed)

### Example: `octopusdeploy_feed` to [`octopusdeploy_nuget_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/nuget_feed)

Assume the following resource:

```terraform
resource "octopusdeploy_feed" "feed" {
  download_attempts              = 1
  download_retry_backoff_seconds = 30
  feed_uri                       = "https://api.nuget.org/v3/index.json"
  feed_type                      = "NuGet"
  is_enhanced_mode               = true
  password                       = "test-password"
  name                           = "Test NuGet Feed (OK to Delete)"
  username                       = "test-username"
}
```

The process to migrate from `v0.7.*` to `v0.8.0` requires three (3) steps:

1. update configuration by replacing the resource (i.e. `octopusdeploy_feed`) with its replacement (i.e. [`octopusdeploy_nuget_feed`](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs/resources/nuget_feed))
2. [`import`](https://www.terraform.io/cli/import) the state to reflect these changes
3. update configuration to reflect your desired changes

Updating the configuration requires a new and empty resource:

```terraform
resource "octopusdeploy_nuget_feed" "feed" {
}
```

Next, import the existing state via the CLI:

```shell
$ terraform import octopusdeploy_nuget_feed.feed "Feeds-123"
```

The address, `octopusdeploy_nuget_feed.feed` will match the resource in your configuration. The ID field, `"Feeds-123"` is the ID of the feed in Octopus Deploy.

Finally, you'll need to update the resource to reflect these changes:

```terraform
resource "octopusdeploy_nuget_feed" "feed" {
  download_attempts              = 1
  download_retry_backoff_seconds = 30
  feed_uri                       = "https://api.nuget.org/v3/index.json"
  is_enhanced_mode               = true
  password                       = "test-password"
  name                           = "Test NuGet Feed (OK to Delete)"
  username                       = "test-username"
}
```