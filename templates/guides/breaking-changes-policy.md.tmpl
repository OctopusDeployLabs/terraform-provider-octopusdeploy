---
page_title: "Breaking Changes Policy"
subcategory: "Upgrades & Migrations"
---

# Our position on Breaking Changes
There are times that we need to make breaking changes to ensure we can maintain the provider and give a great experience for all our Terraform users.

In the context of the Terraform Provider, we define a Breaking Change as one that requires intervention to upgrade from a given version of the provider to another. This intervention might be things like changing your Terraform config, running helper scripts we provide to migrate your State files, or running individual Terraform commands to manually migrate your state files.

When breaking changes do become necessary, we have a process in place to ensure that breaking changes are not surprising to our customers who rely on the Terraform Provider to manage and scale their Octopus usage. This process mirrors the [Octopus Server deprecations process](https://octopus.com/docs/deprecations#overview), implementing the same timeframes and approach. The only difference is that Octopus Server is versioned on a quarterly release cadence, not SemVer - we've made changes for Terraform to support SemVer specifics.

## How we manage Breaking Changes

* We will only make breaking changes when strictly necessary
* We will keep this policy in-sync with the [Octopus Server deprecation process](https://octopus.com/docs/deprecations#overview) for consistency
* There are three key events in the breaking changes timeline:
  * **Announcement**: We will announce planned breaking changes - via Release Notes, documentation and in-provider warnings.
  * **Enactment**: 6 months after announcement, we will enact the deprecation via a "soft-delete". The deprecated features will no longer be available for most users, but an "escape hatch" can temporarily turn them back on to facilitate migration.
  * **Completion**: 12 months after announcement, the deprecation is complete and the features are permanently removed from the code-base.
* We will provide detailed documentation and guides to help you plan the changes you need to make
* We will use [SemVer](https://semver.org/) Major Versions to denote releases that contain breaking changes

## How we'll communicate

* We will keep this documentation up to date with our deprecations process, so you know what to expect
* We will publish Release Notes against releases in GitHub, which will clearly identify any planned and enacted breaking changes. We publish a consolidated summary of these release notes in the Terraform Registry provider documentation under `Upgrades & Migrations`
* We will publish a Migration Guide for each major version in which a Breaking Change is enacted. It will contain details of the steps needed to migrate, and the escape-hatch mechanisms available. These guides are available in the Terraform Registry provider documentation under `Upgrades & Migrations`
* We will ensure warnings and error messages related to Breaking Changes are consistent, descriptive, and provide a link to more detailed information
* We will be available via Octopus Support for specific help

## Version bumps we'll use

* **Announcement**: `Minor` version bump. We will mark the deprecated feature with a warning attribute/message, and where applicable introduce the replacement feature. It's Minor because the deprecation is a warning only, and is a backward-compatible change.
* **Enactment**: `Major` version bump. We will soft-delete the depreacted feature. It's Major because even though it's a soft-delete, majority of customers are expected to perform their migration before upgrading to this version, and we don't want or expect customers to be using the escape-hatch for long.
* **Completion**: `Patch` version bump. We will remove the deprecated feature and its feature flags from the codebase. It's a patch because all customers are expected to have completed their migration in the 12 months since Announcement.

## Example scenario
This is a fictional scenario, but one that demonstrates the process we would go through when making a breaking change
> The `octopusdeploy_tenant` resource has an attribute called `projects`, which contains details of the Projects the Tenant is connected to. To enable more capability in how the provider can be used, we need to extract this `projects` attribute to be a separate resource, and because of details of the underlying API implementation, we can't keep both approaches.
> We plan to introduce a new resource called `octopusdeploy_tenant_project` and remove the existing attribute on the `octopusdeploy_tenant` resource.

### Release and communications timeline
| Date | Event | Customer Can Continue Using Config |
| ---- | ----- | ---------------------------------- |
| 2024-12-09 | **Announcement**<br />Current version of the provider is `v1.2.1`<br /><br />We publish a "Planned Deprecation" for the `octopusdeploy_tenant.projects` attribute, to the Deprecations section of the provider documentation<br /><br />We introduce the new `octopusdeploy_tenant_project` resource, mark the `octopusdeploy_tenant.projects` attribute as Deprecated, and release this deprecation notice as a minor version, `v1.3.0`<br /><br />When they upgrade to `v1.3.0`, customers whose configs use the `octopusdeploy_tenant.projects` attribute receive a warning notifying them that the attribute is deprecated, the date on which it will be removed, and a link to the migration guide<br /><br />The deprecation warning can be suppressed by setting an environment variable named `TF_OCTOPUS_SUPPRESS_DEPRECATION_WARNINGS` to include the string `"octopusdeploy_tenant.projects"`. | ✅ Migration suggested, but not required. |
| ... | Other backwards-compatible features and bug-fixes are added. | ✅ Migration suggested, but not required. |
| 2024-06-09 | **Enactment**<br />Current version of the provider is `v1.6.3`<br /><br />We "soft-delete" the `octopusdeploy_tenant.projects` attribute, and release this soft-deletion as a major version, `v2.0.0`.<br /><br />When they upgrade to `v2.0.0`, customers whose configs use the `octopusdeploy_tenant.projects` now fail to `plan` and `apply`.<br /><br />As an escape-hatch, the soft-deletion can be reversed by setting an environment variable named `TF_OCTOPUS_DEPRECATION_REVERSALS` to include the string `"octopusdeploy_tenant.projects"` (this will be documented in the Deprecation Announcement/Release Notes). This enables customers to continue using their config, but they will still receive the deprecation warning. The deprecation warning suppression environment flag will no longer function for this warning. | ⚠️ Migration required by default, but escape-hatch possible. |
| ... | Other backwards-compatible features and bug-fixes are added. | ⚠️ Migration required by default, but escape-hatch possible. |
| 2025-12-09 | **Completion**<br />Current version of the provider is `v2.1.7`.<br /><br />We remove the `octopusdeploy_tenant.projects` attribute and the escape hatch from the codebase entirely. We publish this change as a patch version, `v2.1.8`.<br /><br />When they upgrade to `v2.1.8`, any customers whose config was still referencing the `octopusdeploy_tenant.projects` attribute will fail to `plan` and `apply` until they perform the migration. | ❌ Migration required. |

