package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func (suite *IntegrationTestSuite) TestGitHubFeed_UpgradeFromSDK_ToPluginFramework() {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testGitHubFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: gitHubconfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   gitHubconfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedGitHubConfig,
				Check: resource.ComposeTestCheckFunc(
					testGitHubFeedUpdated(t),
				),
			},
		},
	})
}

const gitHubconfig = `resource "octopusdeploy_github_repository_feed" "feed_github_repository_migration" {
						  name                                 = "Test GitHub Feed"
						  feed_uri                             = "https://api.github.com"
						  username                             = "username"
						  password                             = "password"
						  download_attempts                    = 6
						  download_retry_backoff_seconds       = 11
					   }`

const updatedGitHubConfig = `resource "octopusdeploy_github_repository_feed" "feed_github_repository_migration" {
						  name                                 = "Updated Test GitHub Feed"
						  feed_uri                             = "https://api.github.com/updated"
						  username                             = "username_Updated"
						  password                             = "password_Updated"
						  download_attempts                    = 7
						  download_retry_backoff_seconds       = 12
					   }`

func testGitHubFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_github_repository_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testGitHubFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_github_repository_feed"+".feed_github_repository_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		githubRepositoryFeed := feed.(*feeds.GitHubRepositoryFeed)

		assert.Regexp(t, "^Feeds\\-\\d+$", githubRepositoryFeed.ID, "Feed ID did not match expected value")
		assert.Equal(t, "Updated Test GitHub Feed", githubRepositoryFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "username_Updated", githubRepositoryFeed.Username, "Feed username did not match expected value")
		assert.Equal(t, true, githubRepositoryFeed.Password.HasValue, "Feed password should be set")
		assert.Equal(t, "https://api.github.com/updated", githubRepositoryFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, 7, githubRepositoryFeed.DownloadAttempts, "Feed download attempts did not match expected value")
		assert.Equal(t, 12, githubRepositoryFeed.DownloadRetryBackoffSeconds, "Feed download retry_backoff_seconds did not match expected value")

		return nil
	}
}
