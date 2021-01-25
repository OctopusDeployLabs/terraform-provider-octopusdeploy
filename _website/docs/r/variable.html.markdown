---
layout: "octopusdeploy"
page_title: "Octopus Deploy: variable"
---

# Resource: Variables

[Variables](https://octopus.com/docs/deployment-process/variables) are values that change based on the
scope of the deployments (e.g. changing SQL Connection Strings between production and staging deployments).

## Example Usage

Basic usage:

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

resource "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  type       = "String"
  value      = "Server=myServerAddress;Database=myDataBase;Trusted_Connection=True;"
}
```

More complex example (with environments and prompts)

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

data "octopusdeploy_environment" "staging" {
    name = "Staging"
}

resource "octopusdeploy_variable" "connection_string" {
  project_id = "${data.octopusdeploy_project.finance.id}"
  name       = "SQLConnectionString"
  type       = "String"
  value      = "Server=myServerAddress;Database=myDataBase;Trusted_Connection=True;"

  scope {
    environments = ["${data.octopusdeploy_environment.staging.id}"]
  }

  prompt{
    label    = "SQL Server Connection String"
    required = false
  }
}
```

Example with sensitive value

```hcl
data "octopusdeploy_project" "finance" {
    name = "Finance"
}

resource "octopusdeploy_variable" "top_secret" {
  project_id       = "${data.octopusdeploy_project.finance.id}"
  name             = "Top Secret Thing"
  type             = "Sensitive"
  is_sensitive     = true
  sensitive_value  = "Server=myServerAddress;Database=myDataBase;Trusted_Connection=True;"
  pgp_key          = "keybase:octopus_user"
}
```

## Argument Reference

* `project_id` (Required) ID of the Project to assign the variable against.
* `name` - (Required) Name of the variable
* `type` - (Required) Type of the variable. Must be one of `String`, `Certificate`, `Sensitive` or `AmazonWebServicesAccount`
* `value` - (Optional) The value of the variable. One of `value` or `sensitive_value` must be set.
* `sensitive_value` - (Optional) The sensitive value of the variable. One of `value` or `sensitive_value` must be set. ~> NOTE: Octopus Deploy server does not return values for Sensitive variables. This means that if the value is changed on the Octopus server Terraform will not be able to detect the drift from your defined configuration.
* `is_sensitive` - (Optional, Default `false`) Whether the variable contains a sensitive value. If this is `true` then `type` must be set to `Sensitive`. 
* `pgp_key` - (Optional) Either a base-64 encoded PGP public key, or a keybase username in the form `keybase:some_person_that_exists`
* `description` - (Optional) Description of the variable
* `scope` - (Optional) The scope to apply to this variable. Contains a list of arrays. All are optional:
    * (Optional) `environments`, `machines`, `actions`, `roles`, `channels`, `tenant_tags`
* `prompt` - (Optional) Prompt for value when a build is run
    * `label` - (Optional) The label for the prompt
    * `description` - (Optional) The description for the prompt
    * `required`- (Optional) Whether or not the value is required

## Attributes reference

* `id` - ID of the variable
* `name` - Name of the variable
* `type` - Type of the variable
* `value` - Value of the variable
* `description` - Description of the variable
* `key_fingerprint` - The fingerprint of the PGP key used to encrypt the secret
* `encrypted_value` - The encrypted value of the secret. ~> NOTE: The encrypted secret may be decrypted using the command line, for example: `terraform output encrypted_value | base64 --decode | keybase pgp decrypt`