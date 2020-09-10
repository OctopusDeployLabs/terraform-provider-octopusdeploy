package main

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: octopusdeploy.Provider})
}
