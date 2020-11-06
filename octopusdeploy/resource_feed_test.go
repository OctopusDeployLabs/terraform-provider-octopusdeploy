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
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployFeed + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	const feedType = "NuGet"
	const feedURI = "http://test.com"
	const enhancedMode = constTrue
	const feedUsername = constUsername
	const feedPassword = constPassword

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testFeedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testFeedtBasic(localName, name, feedType, feedURI, feedUsername, feedPassword, enhancedMode),
				Check: resource.ComposeTestCheckFunc(
					testFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constFeedType, feedType),
					resource.TestCheckResourceAttr(prefix, constFeedURI, feedURI),
					resource.TestCheckResourceAttr(prefix, constUsername, feedUsername),
					resource.TestCheckResourceAttr(prefix, constPassword, feedPassword),
					resource.TestCheckResourceAttr(prefix, constEnhancedMode, enhancedMode),
				),
			},
		},
	})
}

func testFeedtBasic(localName string, name string, feedType string, feedURI string, feedUsername string, feedPassword string, enhancedMode string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name          = "%s"
		feed_type     = "%s"
		feed_uri      = "%s"
		username      = "%s"
		password      = "%s"
		enhanced_mode = "%s"
	}`, constOctopusDeployFeed, localName, name, feedType, feedURI, feedUsername, feedPassword, enhancedMode)
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
