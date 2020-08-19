## Installation
Because the Octopus Deploy Terraform provider is not official and community based, it is considered a third-party provider. Since it's a third-party provider, it must be manually installed.

## Golang Installation
To perform the needed task, you must install Golang. Depending on the operating system you're running, you can simply google **install golang*

# Install the Terraform Provider
The first command you will want to run is `go get` to pull down the executable
1. `go get github.com/OctopusDeploy/terraform-provider-octopusdeploy`

Once the executable is pulled down, it'll automatically go into the `~/go` directory on Linux/MacOS or the `go` directory on the home folder in Windows. Three folders will be shown in the `go` directory:

  - bin
  - pkg
  - src

Typically the executable is in the **bin* directory.

2. `cd` into `~/go/bin`

3. Switching gears for a moment - in the directory where you want the Terraform configuration files to exist to use the Octopus Deploy Terraform provider, create the directory `.terraform/plugins/OS plugin`. The OS Plugin will be different based on OS, so for example, MacOS would look like `.terraform/plugins/darwin_amd64`

4. Copy the `terraform-provider-octopusdeploy` into `.terraform/plugins/OS plugin`

You should know be able to initialize the environment.