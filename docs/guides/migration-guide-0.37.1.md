# v0.37.1 Migration Guide
In this release, we've announced a deprecation that will require action from some customers depending on their configuration

## Deprecated - `octopusdeploy_project.versioning_strategy`
In this release, we announced the deprecation of the `octopusdeploy_project.versioning_strategy` attribute in favour of the newly introduced `octopusdeploy_project_versioning_strategy` resource.

### Rationale
The old attribute was constrained because certain use cases for the the versioning strategy need knowledge of the Deployment Process, which is a separate resource. Trying to configure versioning strategies that reference packages used in a process would cause circular-dependency problems in Terraform that aren't an issue in the Octopus Server Portal. To unlock this scenario, we've extracted the versioning strategy to be its own resource, allowing Terraform to properly plan and apply the dependencies between the Project, the Process and the Versioning Strategy.

### Impact
This change requires some customers to update their HCL. 
Only customers who were already using the `octopusdeploy_project.versioning_strategy` attribute are impacted by this change.

### Timeline
Migration is not yet required, 

| Date       | What we'll do                                                            | What you need to do                                                      |
|------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|
| 2025-06-04 | **Enactment**: Soft-delete the deprecated attribute (Major release)      | Migrate your Terraform config, or use the escape-hatch, before upgrading |
| 2025-12-04 | **Completion**: Remove the deprecated attribute entirely (Patch release) | Migrate your Terraform config before upgrading                           |

### How to migrate
This guide assumes that you already have a Project managed by Terraform, which is using the `versioning_strategy` attribute.

Please ensure you are working from a clean slate and have no pending changes to your Terraform config, by running a `terraform plan`. If you have outstanding changes, please resolve them before proceeding with this guide.

1. Declare a new resource of type `octopusdeploy_project_versioning_strategy`
1. Set the new resource's `project_id` attribute to the existing Project's ID
1. Explicitly set the new resource's `space_id` attribute to the existing Project's Space ID
1. Explicitly set the new resource to depend on the existing Project
1. Transpose the properties from the existing Project's `versioning_strategy` attribute to the new resource. No changes should be required, as we've kept all names and values compatible between the two.
1. Run a `terraform plan`. The only planned changes should be a modification of the Project to remove the `versioning_strategy` attribute, and creation of the new resource.
1. Once you are satisfied with the planned changes, run a `terraform apply` to complete the migration

### Escape hatch

We expect customers to migrate their configs in the 6 months between Announcement and Enactment of a deprecation. However, we know that this isn't always possible, so we have a further 6 months grace period.

If you're caught out during this period and need a bit more time to migrate, you can use this escape hatch to revert the soft-deletion from the Enactment stage.

| Environment Variable | Required Value |
| - | - |
| `TF_OCTOPUS_DEPRECATION_REVERSALS` | `Project-Versioning-Strategy` |

This escape hatch will be removed and migration will be required during the [Completion phase](#Timeline)
