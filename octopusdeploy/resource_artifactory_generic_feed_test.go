package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployArtifactoryGenericFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_artifactory_generic_feed." + localName

	feedURI := "https://example.jfrog.io"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	repository := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	layoutRegex := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testArtifactoryGenericFeedCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testArtifactoryGenericFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
					resource.TestCheckResourceAttr(prefix, "repository", repository),
					resource.TestCheckResourceAttr(prefix, "layout_regex", layoutRegex),
				),
				Config: testArtifactoryGenericFeedBasic(localName, feedURI, name, username, password, repository, layoutRegex),
			},
		},
	})
}

func testArtifactoryGenericFeedBasic(localName string, feedURI string, name string, username string, password string, repository string, layoutRegex string) string {
	return fmt.Sprintf(`resource "octopusdeploy_artifactory_generic_feed" "%s" {
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		username = "%s"
		repository = "%s"
		layout_regex = "%s"
	}`, localName, feedURI, name, password, username, repository, layoutRegex)
}

func testArtifactoryGenericFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testArtifactoryGenericFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_artifactory_generic_feed" {
			continue
		}

		client := testAccProvider.Meta().(*client.Client)
		feed, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Artifactory Generic feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
