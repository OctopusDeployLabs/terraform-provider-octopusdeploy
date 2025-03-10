package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccOctopusDeployMavenFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_maven_feed." + localName

	downloadAttempts := acctest.RandIntRange(1, 10)
	downloadRetryBackoffSeconds := acctest.RandIntRange(0, 60)
	feedURI := "https://repo.maven.apache.org/maven2/"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testMavenFeedCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testMavenFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "download_attempts", strconv.Itoa(downloadAttempts)),
					resource.TestCheckResourceAttr(prefix, "download_retry_backoff_seconds", strconv.Itoa(downloadRetryBackoffSeconds)),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testMavenFeedBasic(localName, downloadAttempts, downloadRetryBackoffSeconds, feedURI, name, username, password),
			},
		},
	})
}

func testMavenFeedBasic(localName string, downloadAttempts int, downloadRetryBackoffSeconds int, feedURI string, name string, username string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_maven_feed" "%s" {
		download_attempts = "%v"
		download_retry_backoff_seconds = "%v"
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		username = "%s"
	}`, localName, downloadAttempts, downloadRetryBackoffSeconds, feedURI, name, password, username)
}

func testMavenFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testMavenFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_maven_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Maven feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
