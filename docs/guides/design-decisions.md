---
page_title: "Design Decisions"
subcategory: "Guides"
---

# Design decisions

This document is to explain our stance on certain implementation decisions for the provider.

## Version-Controlled Projects and Terraform

When managing Deployment Processes and Runbooks you can do two things

1. Define a Database-backed Project and its Process(es) in Terraform
2. Define a Version-Controlled Project and its settings in Terraform, but define the Process as OCL in Git

~> We won’t allow you to try to manage the Process of a Version-Controlled Project in Terraform, and will give an informative error message if you try.

### Why
There should be one source-of-truth for your as-code interactions with Octopus.

Terraform can’t nicely support all the richness that Version Control and Git provide (managing multiple branches, diverging processes etc).

We want to let Version Control do what it does best, and not blur the lines between the two.
