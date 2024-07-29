package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"octopusdeploy": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				octopusdeploy.Provider().GRPCProvider)
			if err != nil {
				log.Fatal(err)
			}

			if err != nil {
				log.Fatal(err)
			}
			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
				providerserver.NewProtocol6(NewOctopusDeployFrameworkProvider()),
			}

			return tf6muxserver.NewMuxServer(context.Background(), providers...)
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	if t.Name() == "TestAccPreCheck" {
		t.Skip("Go registers this function as a test, it's intended as validation")
	}
	if v := os.Getenv("OCTOPUS_URL"); isEmpty(v) {
		t.Fatal("OCTOPUS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("OCTOPUS_APIKEY"); isEmpty(v) {
		t.Fatal("OCTOPUS_APIKEY must be set for acceptance tests")
	}
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
