---
page_title: "Migrating to v0.37.1"
subcategory: "Upgrades & Migrations"
---

# v0.37.1 Migration Guide
In this release, we've announced a deprecation that will require action from some customers, depending on their configuration

## Deprecated - `octopusdeploy_project.versioning_strategy`
In this release, we announced the deprecation of the `octopusdeploy_project.versioning_strategy` attribute in favour of the newly introduced `octopusdeploy_project_versioning_strategy` resource.

### Rationale
The old attribute was constrained because certain use cases for the the versioning strategy need knowledge of the Deployment Process, which is a separate resource. Trying to configure versioning strategies that reference packages used in a process would cause circular-dependency problems in Terraform that aren't an issue in the Octopus Server Portal. To unlock this scenario, we've extracted the versioning strategy to be its own resource, allowing Terraform to properly plan and apply the dependencies between the Project, the Process and the Versioning Strategy.

### Impact
This change requires some customers to update their HCL. 
Only customers who were already using the `octopusdeploy_project.versioning_strategy` attribute are impacted by this change.

### Timeline
Migration will be required no earlier than 2025-12-04

| Date       | What we'll do                                                            | What you need to do                                                      |
|------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|
| 2025-06-04 | **Enactment**: Soft-delete the deprecated attribute (Major release)      | Migrate your Terraform config, or use the escape-hatch, before upgrading |
| 2025-12-04 | **Completion**: Remove the deprecated attribute entirely (Patch release) | Migrate your Terraform config before upgrading                           |

### How to migrate
This guide assumes that you already have a Project managed by Terraform, which is using the `versioning_strategy` attribute.

-> This migration removes and then re-configures the Versioning Strategy on the Project, but this is non-destructive as long as you complete the migration in one go and don't try to create releases between removing the old attribute approach and applying the new resource approach.

~> The original `octopusdeploy_project_versioning_strategy` resource had some incorrect schema definitions that didn't match what Octopus Server API expected, which will make following this upgrade guide impossible if you use only a Release Notes `template`, and not a `donor_package`. We fixed this bug in [Release `v0.40.2`](https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/releases/tag/v0.40.2). If your config requires using the `template` attribute, please ugprade directly to v0.40.2 or higher. There were no breaking changes between v0.37.1 and v0.40.2. 

Please ensure you are working from a clean slate and have no pending changes to your Terraform config, by running a `terraform plan`. If you have outstanding changes, please resolve them before proceeding with this guide.

1. Declare a new resource of type `octopusdeploy_project_versioning_strategy`
1. Set the new resource's `project_id` attribute to the existing Project's ID
1. Explicitly set the new resource's `space_id` attribute to the existing Project's Space ID
1. Explicitly set the new resource to depend on the existing Project
1. Transpose the properties from the existing Project's `versioning_strategy` attribute to the new resource. (See the note above about the `template` attribute if you're using it: you may need to go straight to `v0.40.2` to successfully migrate).
1. Run a `terraform plan`. The only planned changes should be a modification of the Project to remove the `versioning_strategy` attribute, and creation of the new resource.
1. Once you are satisfied with the planned changes, run a `terraform apply` to complete the migration

### Escape hatch

We expect customers to migrate their configs in the 6 months between Announcement and Enactment of a deprecation. However, we know that this isn't always possible, so we have a further 6 months grace period.

If you're caught out during this period and need a bit more time to migrate, you can use this escape hatch to revert the soft-deletion from the Enactment stage.

| Environment Variable | Required Value |
|----------------------|----------------|
| `TF_OCTOPUS_DEPRECATION_REVERSALS` | `Project-Attribute-Versioning-Strategy-0-37-1` |

This escape hatch will be removed and migration will be required during the [Completion phase](#Timeline)
