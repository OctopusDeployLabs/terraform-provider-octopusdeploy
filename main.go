package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/pluginprovider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"log"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	version string = "dev"
	commit  string = "snapshot"
)

func main() {
	//var debugMode bool
	//flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	//flag.Parse()
	//
	//opts := &plugin.ServeOpts{
	//	ProviderFunc: octopusdeploy.Provider,
	//}
	//
	//if debugMode {
	//	opts.Debug = true
	//	opts.ProviderAddr = "octopus.com/com/octopusdeploy"
	//}
	//
	//plugin.Serve(opts)

	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	fullVersion := fmt.Sprintf("%s (%s)", version, commit)

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		octopusdeploy.New(fullVersion)().GRPCProvider,
	)

	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
		providerserver.NewProtocol6(pluginprovider.Provider()),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/OctopusDeployLabs/octopusdeploy",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}

}
