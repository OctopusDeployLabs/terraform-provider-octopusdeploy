package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: octopusdeploy.Provider})
}
