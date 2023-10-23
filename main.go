package main

import (
	"flag"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: octopusdeploy.Provider,
	}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "octopus.com/com/octopusdeploy"
	}

	plugin.Serve(opts)
}
