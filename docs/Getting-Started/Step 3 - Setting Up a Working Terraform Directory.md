# Step 3 - Terraform Directory

Before we can create a resource, we need a place for the HCL code to exist. Of similar importance, we need the Octopus Deploy Terraform provider.

Traditionally when you use Terraform, you can run `terraform init` and the provider gets pulled down for you. However, because the Terraform Provider is not yet in the Hashicorp store, you'll need to manually pull the package down and create the directory that the package should be in.

## Creating the Directory
In the directory where you plan to store the Terraform configuration files, you'll need to create the directory where the provider package will live. 

### On Windows
The directory is:

`.teraform/plugins/windows_amd64`

### On MacOS
The directory is:
`.terraform/plugins/darwin_amd64`

After you create the directory, it should look something like the screenshot below on Mac for example.

![](images/terraformdirectory.jpg)