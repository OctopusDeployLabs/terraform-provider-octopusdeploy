# 🐙 Terraform Provider for Octopus Deploy

## :warning: Warning

The Terraform Provider for Octopus Deploy is under active development and undergoing migration from Terraform SDK to Terraform Plugin Framework. Its functionality can and will change; it is a v0.\* product until its robustness can be assured. Please be aware that types like resources can and will be modified over time. It is strongly recommended to `validate` and `plan` configuration prior to committing changes via `apply`.

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

## Create a New Resource

> [!IMPORTANT]
> We're currently migrating all resources and data sources from Terraform SDK to Terraform Plugin Framework.
> 
> All new resources should be created using Framework, in the `octopusdeploy-framework` directory. [A GitHub action](.github/workflows/prevent-new-sdk-additions.yml) will detect and prevent any new additions to the old `octopusdeploy` SDK directory.

### Acceptance tests
All new resources need an acceptance test that will ensure the lifecycle of the resource works correctly this includes Create, Read, Update and Delete.

### Blocks
> [!WARNING]
> Please avoid using [blocks](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/blocks): these are mainly used for backwards compatability of resources migrated from SDK. Use [nested attributes](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#nested-attribute-types) instead.

### Compatibility with Server
If a resource is not compatible with older versions of the Octopus Deploy server or requires specific feature flags to be enabled, ensure that these requirements are clearly enforced with appropriate validation and descriptive error messaging.

For example to prevent resource usage in versions earlier than 2025.1:
```go
func (f *deploymentFreezeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
  f.Config = ResourceConfiguration(req, resp)

  if f.Config != nil {
    diags := f.Config.EnsureResourceCompatibilityByVersion(deploymentFreezeResourceName, "2025.1")
	resp.Diagnostics.Append(diags...)
  }
}
```

To prevent resource usage based on a feature flag 
```go
func (f *deploymentFreezeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
  f.Config = ResourceConfiguration(req, resp)
	
  if f.Config != nil {
	diags := f.Config.EnsureResourceCompatibilityByFeature(deploymentFreezeResourceName, "ProjectDeploymentFreezesFeatureToggle")
	resp.Diagnostics.Append(diags...)
  }
}
```

## Existing Resource
When modifying an existing SDK resource, we strongly recommend migrating it to Framework first - but this might not always be feasible. We'll judge it on a case-by-case basis.

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

[Running integration tests](running-integration-tests-locally.md)

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
