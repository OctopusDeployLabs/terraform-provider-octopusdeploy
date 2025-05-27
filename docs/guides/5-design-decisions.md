---
page_title: "5. Design Decisions"
subcategory: "Guides"
---

# 5. Creating multiple Octopus Deploy resources

We want explain some of our design decisions to help understand our implementation choices.

## "Configuration as Code" resources
Configuration as Code (config-as-code) is Octopus Deploy feature which allows persist resources in a Git repository. For now, project, deployment process, runbook processes, deployment settings, and non-sensitive variables can be version-controlled.

Resources stored in Git will not be maintained by Terraform Provider and result in "no-op" operations when resource were managed by Provider before being converted to version-controlled  

-> Project is an exception and still can be stored in a Git and managed by Terraform Provider

### Why
We think that Terraform and "Configuration as Code" feature serve same purpose of providing control over Octopus Deploy configuration by storing it as a text in version-controlled system.  
Managing same resource in different Git repositories introduces complex relationship between to sources of truth, which can lead in conflicting dependencies which is difficult to maintain.  

By disabling possibility to manage version-controlled resources in Terraform Provider we explicitly making Configuration as Code a single source of truth for the Octopus Deploy resource.
