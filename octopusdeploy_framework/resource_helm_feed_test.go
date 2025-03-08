package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccOctopusDeployHelmFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_helm_feed." + localName

	feedURI := "https://charts.helm.sh/stable"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testHelmFeedCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testHelmFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testHelmFeedBasic(localName, feedURI, name, username, password),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testHelmFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", updatedName),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testHelmFeedBasic(localName, feedURI, updatedName, username, password),
			},
		},
	})
}

func testHelmFeedBasic(localName string, feedURI string, name string, username string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_helm_feed" "%s" {
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		username = "%s"
	}`, localName, feedURI, name, password, username)
}

func testHelmFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testHelmFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_helm_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Helm feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
