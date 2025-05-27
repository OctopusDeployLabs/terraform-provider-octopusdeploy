---
page_title: "5. Design Decisions"
subcategory: "Guides"
---

# 5. Design decisions

We want explain some of our design decisions to help understand our implementation choices.

## "Configuration as Code" resources

### Overview
The Octopus Deploy Terraform provider has been designed with careful consideration for how it interacts with Octopus Deploy's native "Configuration as Code" (Config as Code) feature. This section explains our approach to managing resources that can exist in both systems.

### How Config as Code Works with Terraform
Configuration as Code in Octopus Deploy allows specific resources to be stored in a Git repository rather than in the Octopus database. Currently, this includes projects, deployment and runbook processes, deployment settings, and non-sensitive variables.

### When a resource is version-controlled via Config as Code
* Terraform operations become "no-op" for that resource: If a resource was previously managed by Terraform and later converted to be version-controlled, Terraform operations on that resource will have no effect.

* Projects are an exception: Projects can be both version-controlled and managed by Terraform simultaneously. This allows you to configure project settings via Terraform while keeping deployment processes in Git.

### Why This Approach?
We've adopted this design pattern for several important reasons:

* Avoiding Conflicting Sources of Truth: Both Terraform and Config as Code serve the purpose of maintaining infrastructure in a version-controlled text format. Having the same resource managed in different Git repositories would create competing sources of truth.

* Simplified Mental Model: By clearly delineating which system controls which resources, we create a clearer mental model for practitioners using both systems.

* Preventing Configuration Drift: Attempting to manage the same resources in both systems would inevitably lead to drift and reconciliation challenges.

### Recommended Practices
For the best experience using Terraform with Octopus Deploy:

* Use Terraform for infrastructure setup, project creation, and global configuration
* Use Configuration as Code for deployment and runbook processes and variables
* For projects that use Config as Code, use Terraform only for the project settings, not for the elements stored in Git
* Be aware that converting a Terraform-managed resource to use Config as Code will make Terraform operations on that resource ineffective
