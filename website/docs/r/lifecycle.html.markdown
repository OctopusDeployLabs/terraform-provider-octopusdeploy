---
layout: "octopusdeploy"
page_title: "Octopus Deploy: lifecycle"
---

# Resource: octopusdeploy_lifecycle

Use this resource allows the creation of Octopus Deploy [lifecycles](https://octopus.com/docs/deployment-process/lifecycles).

Lifecycles can be used to automatically promote deployments between environments, and limit environments that can be deployed to until a release has been thoroughly tested.

## Example Usage

```hcl
resource "octopusdeploy_environment" "Env1" {
  name = "LifecycleTestEnv1"
}

resource "octopusdeploy_environment" "Env2" {
  name = "LifecycleTestEnv2"
}

resource "octopusdeploy_environment" "Env3" {
  name = "LifecycleTestEnv3"
}

resource "octopusdeploy_lifecycle" "foo" {
  name        = "Funky Lifecycle"
  description = "Funky Lifecycle description"

  release_retention_policy {
    unit             = "Items"
    quantity_to_keep = 2
  }

  tentacle_retention_policy {
    unit             = "Days"
    quantity_to_keep = 1
  }

  phase {
    name                                  = "P1"
    minimum_environments_before_promotion = 2
    is_optional_phase                     = true
    automatic_deployment_targets          = ["${octopusdeploy_environment.Env1.id}"]
    optional_deployment_targets           = ["${octopusdeploy_environment.Env2.id}"]
  }

  phase {
    name = "P2"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the lifecycle.

* `description` - (Optional) Description of the lifecycle.

* `release_retention_policy` - (Optional) A release retention policy block as documented below.

* `phase` - (Optional) A phase block as documented below.

Release Retention Policy (`release_retention_policy`) blocks support the following:

* `unit` - (Optional) The unit of quantity to keep. Either `Days` or `Items`. Defaults to `Days`.

* `quantity_to_keep` - (Optional) The number of units required before a release can enter the next phase. If 0, all environments are required. Defaults to `0`.

Phase (`phase`) blocks support the following:

* `name` - (Required) The name of the phase.

* `minimum_environments_before_promotion` - (Optional) The number of days/releases to keep. If 0 all are kept. Defaults to `0`.

* `is_optional_phase` - (Optional) If false a release must be deployed to this phase before it can be deployed to the next phase. Defaults to `false`.

* `automatic_deployment_targets` - (Optional) Environment Ids in this phase that a release is automatically deployed to when it is eligible for this phase.

* `optional_deployment_targets` - (Optional) Environment Ids in this phase that a release can be deployed to, but is not automatically deployed to.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the environment.
