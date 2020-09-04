# Create an account with Terraform

The **Account Creation** directory contains the following files:

  - main.tf
  - terraform.tfvars
  - variables.tf

## Files

The Terraform configuration files are used to build a specific aspect of the environment. Below is an explanation of the Configuration Files.

### main.tf

The `main.tf` configuration file has two blocks, one for the provider and another for the `account` resource. The provider calls for three inputs:

  address = The server URL
  apikey  = The API key to authenticate to an Octopus Deploy server
  space   = The Octopus Deploy space to create the resource in

The resource itself that's being created is done with the `octopusdeploy_account` Terraform resource. The values in the input blocks are:

  - name            = 
  - account_type    = 
  - subscription_id = 

### variables.tf

The `variables.tf` configuration file contains the following variables to be used in the `main.tf` configuration file.

  - apiKey: The API key to authenticate to an Octopus Deploy server.
  - space: The Octopus Deploy space to create the resource in.
  - serverURL: The server URL.
  - azureAccountName: The metadata name for the Azure account that is being added.
  - subID: The subscription ID for the Azure account.
  - clientID: The client ID that's being used to authenticate to Azure from Octopus Deploy.
  - clientSecret: The client secret of the client ID

### terraform.tfvars

The `terraform.tfvars` configure file contains the values for the variables that are being passed in at runtime. An example is below:

```
serverURL = "https://mikelevanoctopus.octopus.app"
space     = "Default"
azureAccountName = "AzureAC"
subID = "f31cfbfd-37fb-4e3d-a713-7bd6ceb2f7bf"
clientID = "12345-6789-1876-8878-ce8282763205"
clientSecret = "LQaZkf7Vv
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