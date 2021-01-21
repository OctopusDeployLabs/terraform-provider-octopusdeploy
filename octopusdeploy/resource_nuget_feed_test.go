package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployNuGetFeedBasic(t *testing.T) {
	const feedPrefix = "octopusdeploy_nuget_feed.foo"
	const feedName = "Testing Nuget one two three"
	const feedURI = "http://test.com"
	const enhancedMode = true
	const feedUsername = "username"
	const feedPassword = "password"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployNuGetFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testNuGetFeedBasic(feedName, feedURI, feedUsername, feedPassword, enhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployNuGetFeedExists(feedPrefix),
					resource.TestCheckResourceAttr(feedPrefix, "name", feedName),
					resource.TestCheckResourceAttr(feedPrefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(feedPrefix, "username", feedUsername),
					resource.TestCheckResourceAttr(feedPrefix, "password", feedPassword),
					resource.TestCheckResourceAttr(feedPrefix, "is_enhanced_mode", strconv.FormatBool(enhancedMode)),
				),
			},
		},
	})
}

func testNuGetFeedBasic(name, feedURI string, username string, password string, isEnhancedMode bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_nuget_feed" "foo" {
		feed_uri = "%s"
		is_enhanced_mode = %v
		name = "%s"
		password = "%s"
		username = "%s"
	}`, feedURI, isEnhancedMode, name, password, username)
}

func testOctopusDeployNuGetFeedExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		feedID := s.RootModule().Resources[n].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testOctopusDeployNuGetFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_nuget_feed" {
			continue
		}

		_, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
