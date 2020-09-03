# Create a variable with Terraform

The **variable Creation** directory contains the following files:

  - main.tf
  - terraform.tfvars
  - variables.tf

## Files

The Terraform configuration files are used to build a specific aspect of the environment. Below is an explanation of the Configuration Files.

### main.tf

The `main.tf` configuration file has two blocks, one for the provider and another for the `variable` resource. The provider calls for three inputs:

  address = The server URL
  apikey  = The API key to authenticate to an Octopus Deploy server
  space   = The Octopus Deploy space to create the resource in

The resource itself that's being created is done with the `octopusdeploy_variable` Terraform resource. The values in the input blocks are:

  - name       =
  - project_id =
  - type       = "String"

### variables.tf

The `variables.tf` configuration file contains the following variables to be used in the `main.tf` configuration file.

  - apiKey: The API key to authenticate to an Octopus Deploy server.
  - space: The Octopus Deploy space to create the resource in.
  - serverURL: The server URL.
  - projectID: ID of the project that the variable will be created in.
  - varName: Variable name

### terraform.tfvars

The `terraform.tfvars` configure file contains the values for the variables that are being passed in at runtime. An example is below:

```
serverURL   = "https://mikelevanoctopus.octopus.app"
space       = "Default"
projectID = "Projects-1"
varName   = "DevKey"
```

## Using the Project Example

To use Terraform to create the resource:

1. Clone the repository
2. Change directory (`cd`) into the *examples* directory
3. `cd` into the appropriate resource 
4. Edit the `terraform.tfvars` configuration file for the appropriate values
5. Run the following commands to initiate, plan, and apply the creation:
   - `terraform init`
   - `terraform plan`
   - `terraform apply`