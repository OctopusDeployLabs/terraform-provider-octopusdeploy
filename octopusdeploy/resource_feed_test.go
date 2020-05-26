package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployFeedBasic(t *testing.T) {
	const feedPrefix = "octopusdeploy_feed.foo"
	const feedName = "Testing one two three"
	const feedType = "NuGet"
	const feedURI = "http://test.com"
	const enhancedMode = "true"
	const feedUsername = "username"
	const feedPassword = "password"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testFeedtBasic(feedName, feedType, feedURI, feedUsername, feedPassword, enhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployFeedExists(feedPrefix),
					resource.TestCheckResourceAttr(
						feedPrefix, "name", feedName),
					resource.TestCheckResourceAttr(
						feedPrefix, "feed_type", feedType),
					resource.TestCheckResourceAttr(
						feedPrefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(
						feedPrefix, "username", feedUsername),
					resource.TestCheckResourceAttr(
						feedPrefix, "password", feedPassword),
					resource.TestCheckResourceAttr(
						feedPrefix, "enhanced_mode", enhancedMode),
				),
			},
		},
	})
}

func testFeedtBasic(name, feedType, feedURI string, feedUsername string, feedPassword string, enhancedMode string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_feed" "foo" {
			name          = "%s"
			feed_type     = "%s"
			feed_uri      = "%s"
			username = "%s"
			password = "%s"
			enhanced_mode = "%s"
		}
		`,
		name, feedType, feedURI, feedUsername, feedPassword, enhancedMode,
	)
}

func testOctopusDeployFeedExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return feedExistsHelper(s, client)
	}
}

func testOctopusDeployFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyFeedHelper(s, client)
}
