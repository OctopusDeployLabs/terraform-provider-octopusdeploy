package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"octopusdeploy": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OCTOPUS_URL"); isEmpty(v) {
		t.Fatal("OCTOPUS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("OCTOPUS_APIKEY"); isEmpty(v) {
		t.Fatal("OCTOPUS_APIKEY must be set for acceptance tests")
	}
}

func ProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"octopusdeploy": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				Provider().GRPCProvider)
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
				providerserver.NewProtocol6(octopusdeployv6.NewOctopusDeployProviderV6()),
			}

			return tf6muxserver.NewMuxServer(context.Background(), providers...)
		},
	}
}
