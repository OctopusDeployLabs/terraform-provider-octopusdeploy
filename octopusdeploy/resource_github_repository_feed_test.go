package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestGitHubRepositoryFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_github_repository_feed." + localName

	downloadAttempts := acctest.RandIntRange(0, 10)
	downloadRetryBackoffSeconds := acctest.RandIntRange(0, 60)
	feedURI := "https://api.github.com"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testGitHubRepositoryFeedCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testGitHubRepositoryFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "download_attempts", strconv.Itoa(downloadAttempts)),
					resource.TestCheckResourceAttr(prefix, "download_retry_backoff_seconds", strconv.Itoa(downloadRetryBackoffSeconds)),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testGitHubRepositoryFeedBasic(localName, downloadAttempts, downloadRetryBackoffSeconds, feedURI, name, username, password),
			},
		},
	})
}

func testGitHubRepositoryFeedBasic(localName string, downloadAttempts int, downloadRetryBackoffSeconds int, feedURI string, name string, username string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_github_repository_feed" "%s" {
		download_attempts = "%v"
		download_retry_backoff_seconds = "%v"
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		username = "%s"
	}`, localName, downloadAttempts, downloadRetryBackoffSeconds, feedURI, name, password, username)
}

func testGitHubRepositoryFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testGitHubRepositoryFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_github_repository_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("GitHub repository feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestGithubFeedResource verifies that a nuget feed can be reimported with the correct settings
func TestGithubFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "44-githubfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "44a-githubfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := feeds.FeedsQuery{
			PartialName: "Github",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Github\"")
		}
		resource := resources.Items[0].(*feeds.GitHubRepositoryFeed)

		if resource.FeedType != "GitHub" {
			t.Fatal("The feed must have a type of \"GitHub\"")
		}

		if resource.Username != "test-username" {
			t.Fatal("The feed must have a username of \"test-username\"")
		}

		if resource.DownloadAttempts != 1 {
			t.Fatal("The feed must be have a downloads attempts set to \"1\"")
		}

		if resource.DownloadRetryBackoffSeconds != 30 {
			t.Fatal("The feed must be have a downloads retry backoff set to \"30\"")
		}

		if resource.FeedURI != "https://api.github.com" {
			t.Fatal("The feed must be have a feed uri of \"https://api.github.com\"")
		}

		foundExecutionTarget := false
		foundServer := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "Server" {
				foundServer = true
			}
		}

		if !(foundExecutionTarget && foundServer) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"Server\"")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "44a-githubfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
