package main

import (
	"context"
	"flag"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework"
	"log"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		octopusdeploy.Provider().GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(octopusdeploy_framework.NewOctopusDeployFrameworkProvider()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	opts := []tf6server.ServeOpt{}

	var providerName = "registry.terraform.io/OctopusDeployLabs/octopusdeploy"
	if debugMode {
		opts = append(opts, tf6server.WithManagedDebug())
		providerName = "octopus.com/com/octopusdeploy"
	}

	err = tf6server.Serve(providerName, muxServer.ProviderServer, opts...)
	if err != nil {
		log.Fatal(err)
	}
}
