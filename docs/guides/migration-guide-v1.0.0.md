---
page_title: "Migrating to v1.0.0"
subcategory: "Upgrades & Migrations"
---

# v1.0.0 Migration Guide
In this release, we've announced a deprecation that will require action from some customers, depending on their configuration

## Deprecated - `octopusdeploy_deployment_process`
In this release, we announced the deprecation of the `octopusdeploy_deployment_process` resource in favour of the newly introduced `octopusdeploy_process` resource.

## Deprecated - `octopusdeploy_runbook_process`
In this release, we announced the deprecation of the `octopusdeploy_runbook_process` resource in favour of the newly introduced `octopusdeploy_process` resource.

### Rationale
The old resources are prone to state drift issues due to how the schema was defined within in the provider, this would often lead to issues such as not being able to reorder or add additional steps to the process. In addition, the old resources design did not allow for new features such as the ability to use step templates easily.

### Impact
This change requires some customers to update their HCL.
Only customers who were already using the `octopusdeploy_deployment_process` or `octopusdeploy_runbook_process` resources are impacted by this change.

### Timeline
Migration will be required no earlier than 2026-05-10

| Date       | What we'll do                                                            | What you need to do                                                      |
|------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|
| 2025-11-05 | **Enactment**: Soft-delete the deprecated attribute (Major release)      | Migrate your Terraform config, or use the escape-hatch, before upgrading |
| 2026-05-05 | **Completion**: Remove the deprecated resources entirely (Patch release) | Migrate your Terraform config before upgrading                           |

### How to migrate

Please ensure you are working from a clean slate and have no pending changes to your Terraform config, by running a `terraform plan`. If you have outstanding changes, please resolve them before proceeding with this guide.

-> This migration removes the old deployment process and replaces it with the process resource, this is non-destructive as long as you complete the migration in one go.

1. Declare a new resource of type `octopusdeploy_process`
1. Set the `project_id` to the existing Project's ID and if migrating a runbook process set the `runbook_id` to the existing Runbook's ID
1. Declare a new resource of type `octopusdeploy_process_step`
1. Set the `process_id` to the `id` of the new process resource
1. Transpose the `name`
1. Transpose the `type`, `type` is the type of embedded action which corresponds to `action_type` in the `action` block, for built-in action like `run_script_action` this is hidden.
1. Transpose the `properties` and `execution_properties` to the new resource from a `step` and `action` in the old deployment process, `properties` are those that are set on the step level and `execution_properties` are those historically set on the `action`. Note that the process step name of the embedded action is now always the same as the name of the step, this can introduce some changes during migration when the name of the first action of a step is different from its step name.
1. Transpose `run_on_server` to the `execution_properties` using the key `Octopus.Action.RunOnServer`
1. Transpose `window_size` to the `execution_properties` using the key `Octopus.Action.TargetRoles`
1. Transpose `target_roles` to the `properties` using the key `Octopus.Action.TargetRoles`
1. Transpose `primary_package` to the `packages` attribute with the key being an empty string `""`
1. Repeat until all parent `steps` have been transposed
1. For child steps create a new resource of type `octopusdeploy_process_child_step` (nested actions)
1. Set the `process_id` to the `id` of the new process resource
1. Set the `parent_id` to the owning `octopusdeploy_process_step.id`
1. Repeat until all child `steps` (actions) have been transposed
1. Declare a new `octopusdeploy_process_steps_order` resource
1. Set the `process_id` to the `id` of the new process resource
1. Set the list of `steps` to a list of `octopusdeploy_process_step.id` in the order you would like the steps to execute
1. For child `step` order declare a new `octopusdeploy_process_child_steps_order`, do this for each group of child steps.
1. Set the `process_id` to the `id` of the new process resource
1. Set the `parent_id` to the owning `octopusdeploy_process_step.id`
1. Set the list of `children` to a list of `octopusdeploy_process_child_step.id` in the order you would like the child steps to execute
1. Note down the existing `process_id`, `step_id`'s and `action_id`'s stored in the `tfstate` file deployment process being migrated
1. Run `terraform state rm [options] 'octopusdeploy_deployment_process.<name>'` to remove the old deployment process from state
1. Remove the old deployment process from the terraform config
1. Run `terraform import [options] octopusdeploy_process.<name> <process-id>` to load the new process into state
1. Run `terraform import [options] octopusdeploy_process_step.<name> "<process-id>:<step-id>"` for each step to load them into state
1. Run `terraform import [options] octopusdeploy_process_steps_order.<name> <process-id>` to load the process step oder into state
1. Run `terraform import [options] octopusdeploy_process_child_step.<name> "<process-id>:<parent-step-id>:<child-step-id>"` for each child step to load them into order
1. Run `terraform import [options] octopusdeploy_process_child_steps_order.<name> "<process-id>:<parent-step-id>"` for each child step group to load them into state
1. Run `terraform plan` to see if anything is missing from state, fix accordingly
1. When satisfied run `terraform apply` to complete the migration

### Escape hatch

We expect customers to migrate their configs in the 6 months between Announcement and Enactment of a deprecation. However, we know that this isn't always possible, so we have a further 6 months grace period.

If you're caught out during this period and need a bit more time to migrate, you can use this escape hatch to revert the soft-deletion from the Enactment stage.

| Environment Variable | Required Value                 |
|----------------------|--------------------------------|
| `TF_OCTOPUS_DEPRECATION_REVERSALS` | `Process_v1.0.0` |

This escape hatch will be removed and migration will be required during the [Completion phase](#Timeline)
