package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployFeedBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_feed." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	const feedType = "NuGet"
	const feedURI = "http://test.com"
	const enhancedMode = "true"
	const feedUsername = "username"
	const feedPassword = "password"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testFeedtBasic(localName, name, feedType, feedURI, feedUsername, feedPassword, enhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "feed_type", feedType),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "username", feedUsername),
					resource.TestCheckResourceAttr(prefix, "password", feedPassword),
					resource.TestCheckResourceAttr(prefix, "is_enhanced_mode", enhancedMode),
				),
			},
		},
	})
}

func testFeedtBasic(localName string, name string, feedType string, feedURI string, feedUsername string, feedPassword string, enhancedMode string) string {
	return fmt.Sprintf(`resource "octopusdeploy_feed" "%s" {
		name             = "%s"
		feed_type        = "%s"
		feed_uri         = "%s"
		username         = "%s"
		password         = "%s"
		is_enhanced_mode = "%s"
	}`, localName, name, feedType, feedURI, feedUsername, feedPassword, enhancedMode)
}

func testFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		feedID := rs.Primary.ID
		feed, err := client.Feeds.GetByID(feedID)
		if err == nil {
			if feed != nil {
				return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
