package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

type ociRegistryFeedTestData struct {
	name     string
	uri      string
	username string
	password string
}

func TestAccOctopusDeployOCIRegistryFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_oci_registry_feed." + localName
	data := ociRegistryFeedTestData{
		name:     acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		uri:      "oci://integration-test-registry.docker.io",
		username: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		password: acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testOCIRegistryFeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testOCIRegistryFeedBasic(data, localName),
				Check:  testAssertOCIRegistryAttributes(data, prefix),
			},
			{
				Config: testOCIRegistryFeedBasicUpdate(data, localName),
				Check:  testAssertOCIRegistryAttributes(data, prefix),
			},
		},
	})
}

func testAssertOCIRegistryAttributes(expected ociRegistryFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "feed_uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "username", expected.username),
		resource.TestCheckResourceAttr(prefix, "password", expected.password),
	)
}

func testOCIRegistryFeedBasic(data ociRegistryFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_oci_registry_feed" "%s" {
			name		= "%s"
			feed_uri	= "%s"
			username	= "%s"
			password	= "%s"
		}
	`,
		localName,
		data.name,
		data.uri,
		data.username,
		data.password,
	)
}

func testOCIRegistryFeedBasicUpdate(data ociRegistryFeedTestData, localName string) string {
	data.name = data.name + "-updated"

	return testOCIRegistryFeedBasic(data, localName)
}

func testOCIRegistryFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_oci_registry_feed" {
			continue
		}

		feed, err := feeds.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("OCI Registry feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
