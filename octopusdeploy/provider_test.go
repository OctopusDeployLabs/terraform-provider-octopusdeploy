package octopusdeploy

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]schema.Provider{
		"octopusdeploy": *testAccProvider,
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
