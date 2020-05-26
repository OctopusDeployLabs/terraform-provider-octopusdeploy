package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployNugetFeedBasic(t *testing.T) {
	const feedPrefix = "octopusdeploy_nuget_feed.foo"
	const feedName = "Testing Nuget one two three"
	const feedURI = "http://test.com"
	const enhancedMode = "true"
	const feedUsername = "username"
	const feedPassword = "password"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployNugetFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNugetFeedBasic(feedName, feedURI, feedUsername, feedPassword, enhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployNugetFeedExists(feedPrefix),
					resource.TestCheckResourceAttr(
						feedPrefix, "name", feedName),
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

func testNugetFeedBasic(name, feedURI string, feedUsername string, feedPassword string, enhancedMode string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_nuget_feed" "foo" {
			name          = "%s"
			feed_uri      = "%s"
			username = "%s"
			password = "%s"
			enhanced_mode = "%s"
		}
		`,
		name, feedURI, feedUsername, feedPassword, enhancedMode,
	)
}

func testOctopusDeployNugetFeedExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return feedExistsHelper(s, client)
	}
}

func testOctopusDeployNugetFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyFeedHelper(s, client)
}
