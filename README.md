# 🐙 Terraform Provider for Octopus Deploy

## :warning: Warning

The Terraform Provider for Octopus Deploy is under active development. Its functionalty can and will change; it is a v0.\* product until its robustness can be assured. Please be aware that types like resources can and will be modified over time. It is strongly recommended to `validate` and `plan` configuration prior to committing changes via `apply`.

## About

This repository contains the source code for the Terraform Provider for [Octopus Deploy](https://octopus.com). It supports provisioning/configuring of Octopus Deploy instances via [Terraform](https://www.terraform.io/). Documentation and guides for using this provider are located on the Terraform Registry: [Documentation](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs).

## 🪄 Installation and Configuration

The Terraform Provider for Octopus Deploy is available via the Terraform Registry: [OctopusDeployLabs/octopusdeploy](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy). To install this provider, copy and paste this code into your Terraform configuration:

```hcl
terraform {
  required_providers {
    octopusdeploy = {
      source = "OctopusDeployLabs/octopusdeploy"
      version = "version-number" # example: 0.21.1
    }
  }
}

provider "octopusdeploy" {
  # configuration options
  address  = "https://octopus.example.com"     # (required; string) the service endpoint of the Octopus REST API
  api_key  = "API-XXXXXXXXXXXXX"               # (required; string) the API key to use with the Octopus REST API
  space_id = "Spaces-1"                        # (optional; string) the space ID in Octopus Deploy
}
```

If `space_id` is not specified the Terraform Provider for Octopus Deploy will assume the default space.

### Environment Variables

You can provide your Octopus Server configuration via the `OCTOPUS_URL` and `OCTOPUS_APIKEY` environment variables, representing your Octopus Server address and API Key, respectively.

```hcl
provider "octopusdeploy" {}
```

Run `terraform init` to initialize this provider and enable resource management.

## 🛠 Build Instructions

A build of this Terraform Provider can be created using the [Makefile](https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/blob/master/Makefile) provided in the source:

```shell
$ make
```

This will generate a binary that will be installed to the local plugins folder. Once installed, the provider may be used through the following configuration:

```hcl
terraform {
  required_providers {
    octopusdeploy = {
      source  = "octopus.com/com/octopusdeploy"
      version = "0.7.63"
    }
  }
}

provider "octopusdeploy" {
  address  = # address
  api_key  = # API key
  space_id = # space ID
}
```

After the provider has been built and saved to the local plugins folder, it may be used after initialization:

```shell
$ terraform init
```

Terraform will scan the local plugins folder directory structure (first) to qualify the source name provided in your Terraform configuration. If it can resolve this name then the local copy will be initialized for use with Terraform. Otherwise, it will scan the Terraform Registry.

:warning: The `version` number specified in your Terraform configuration MUST match the version number specified in the Makefile. Futhermore, this version MUST either be incremented for each local re-build; otherwise, Terraform will use the cached version of the provider in the `.terraform` folder. Alternatively, you can simply delete the folder and re-run the `terraform init` command.

## Debugging 
If you want to debug the provider follow these steps!

### Prerequisites
- Terraform provider is configured to use the local version e.g. `"octopus.com/com/octopusdeploy"`
```hcl
terraform {
  required_providers {
    octopusdeploy = {
      source  = "octopus.com/com/octopusdeploy"
      version = "0.7.63"
    }
  }
}
```
- Optional - Install delve https://github.com/go-delve/delve

### Via Goland
1. Debug the provided run configuration `Run provider` - This will run the provider with the `-debug` flag set to true.
2. Export the environment variable that the running provider logs out, it will look something like this:
```shell
TF_REATTACH_PROVIDERS='{"octopus.com/com/octopusdeploy":{"Protocol":"grpc","ProtocolVersion":5,"Pid":37485,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/pq/_cv_xzg97ts8t2tq25d_43wr0000gn/T/plugin447612806"}}}'
```
3. In the same terminal session where you exported the environment variable, execute the Terraform commands you want to debug.

### Via Delve
1. Add your breakpoints, this can be done by adding `runtime.Breakpoint()` lines to where you want the code to break.
2. Run `dlv debug . -- --debug` in the root folder of the project (same directory where `main.go` lives).
3. The debugger will start and wait, type `continue` in the terminal to get it to start the provider.
4. Export the environment variable that the running provider logs out, it will look something like this:
```shell
TF_REATTACH_PROVIDERS='{"octopus.com/com/octopusdeploy":{"Protocol":"grpc","ProtocolVersion":5,"Pid":37485,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/pq/_cv_xzg97ts8t2tq25d_43wr0000gn/T/plugin447612806"}}}'
```
5. In the same terminal session where you exported the environment variable, execute the Terraform commands you want to debug.

## Testing

The tests under `integration_test.go` verify the Terraform modules under the `terraform` directory run as expected. At a high level, the tests work like this:

* Create a new Octopus Deploy instance and MSSQL in Docker using Test Containers.
* Create a new blank space using the module in `1-singlespace`.
* Populate the new space using one of the other modules under the `terraform` directory.
* Inspect the new space with the Octopus Go client to verify the newly created resources exist and have the expected settings.

These tests are split by [go-test-split-action](https://github.com/hashicorp-forge/go-test-split-action) to run in parallel across multiple GitHub Action workers as part of the build.

To run the tests locally, you must have the following environment variables:

* `LICENSE`: base 64 encoded octopus license XML file
* `ECR_ACCESS_KEY`: aws access key (optional, and only used for ECR feed tests)
* `ECR_SECRET_KEY`: aws secret key (optional, and only used for ECR feed tests)
* `GIT_CREDENTIAL`: github token (optional, and only used for CaC tests)
* `GIT_USERNAME`: github username (optional, and only used for CaC tests)

By default, the tests spin up an Octopus instance based on the locally pulled `octopusdeploy/octopusdeploy:latest` image. You can override the image and tag with the following environment variables:

* `OCTOTESTIMAGEURL`: The image to use to create the Octopus instance. Defaults to `octopusdeploy/octopusdeploy`. Can be set to `octopusdeploy/linuxoctopus` to test the Octopus cloud builds. Note you must be logged into DockerHub with the correct credentials to pull this private image.
* `OCTOTESTVERSION`: The version of the image to use. Defaults to `latest`.

The tests can be run in parallel to speed up the test run time. Set the `GOMAXPROCS` environment variable to the number of parallel tests you want to run. For example, to run 4 tests in parallel, set `GOMAXPROCS=4`. Note though that each test creates it own Octopus and MSSQL docker images, so you will need to have enough resources to run the tests in parallel.

The tests work by executing the `terraform` executable against the test modules. By default, the version of the terraform provider that is used for these tests is defined in the file `config.tf`. The version defined in this file is almost always going to be a few revisions old.

To test the latest build of the Terraform provider, you must build the provider executable locally and define an override to configure Terraform to ignore the version of the provider defined in the module and instead use your local build.

Build the terraform provider with the command:

```bash
go build -o terraform-provider-octopusdeploy main.go
```

Then save the following configuration to a file called `~/.terraformrc`, making sure to replace the string `/var/home/yourname/Code/` with the path to the directory containing the Terraform provider executable:

```hcl
provider_installation {
  dev_overrides {
    "octopusdeploylabs/octopusdeploy" = "/var/home/yourname/Code/terraform-provider-octopusdeploy"
  }

  direct {}
}
```

When the overrides are in effect, `terraform` will print a warning message to the console:

```bash
Warning: Provider development overrides are in effect
        
The following provider development overrides are set in the CLI
configuration:
 - octopusdeploylabs/octopusdeploy in /var/home/yourname/Code/terraform-provider-octopusdeploy
```

If you wish to view the Octopus instance created by the test framework, add a breakpoint after the call to `testFramework.Act()` and get the IP address and port from the `container` argument:

![](test-debugging.png)

You can then log into the Octopus instance with the credentials `admin` and `Password01!`.

## Documentation Generation

Documentation is auto-generated by the [tfplugindocs CLI](https://github.com/hashicorp/terraform-plugin-docs). To generate the documentation, run the following command:

```shell
$ make docs
```

or
```shell
go generate main.go
```

## 🤝 Contributions

Contributions are welcome! :heart: Please read our [Contributing Guide](CONTRIBUTING.md) for information about how to get involved in this project.
