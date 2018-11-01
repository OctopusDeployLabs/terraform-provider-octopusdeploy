package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/pawelpabich/terraform-provider-octopusdeploy/octopusdeploy"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: octopusdeploy.Provider})
}
