package octopusdeployv6

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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var ()
var testAccProvider provider.Provider
var testAccProviderFactories map[string]func() (provider.Provider, diag.Diagnostics)

func checkEnvVar(t *testing.T, key string) {
	if v := os.Getenv(key); v == "" {
		t.Fatalf("%s must be set for acceptance tests", key)
	}
}

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
				providerserver.NewProtocol6(NewOctopusDeployProviderV6()),
			}

			return tf6muxserver.NewMuxServer(context.Background(), providers...)
		},
	}
}

func TestAccPreCheck(t *testing.T) {
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
