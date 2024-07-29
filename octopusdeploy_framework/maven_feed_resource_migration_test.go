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

func (suite *IntegrationTestSuite) TestMavenResource_UpgradeFromSDK_ToPluginFramework() {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: mavenConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   mavenConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedMavenConfig,
				Check: resource.ComposeTestCheckFunc(
					testFeedUpdated(t),
				),
			},
		},
	})
}

const mavenConfig = `resource "octopusdeploy_maven_feed" "feed_maven_migration" {
						  name                                 = "Maven"
						  feed_uri                             = "https://repo.maven.apache.org/maven2/"
						  username                             = "username"
						  password                             = "password"
						  download_attempts                    = 6
						  download_retry_backoff_seconds       = 11
					   }`

const updatedMavenConfig = `resource "octopusdeploy_maven_feed" "feed_maven_migration" {
						  name                                 = "Updated_Maven"
						  feed_uri                             = "https://Updated.maven.apache.org/maven2/z"
						  username                             = "username_Updated"
						  password                             = "password_Updated"
						  download_attempts                    = 7
						  download_retry_backoff_seconds       = 12
					   }`

func testFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_maven_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_maven_feed"+".feed_maven_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		mavenFeed := feed.(*feeds.MavenFeed)

		assert.Equal(t, "Feeds-1001", mavenFeed.ID, "Feed ID did not match expected value")
		assert.Equal(t, "Updated_Maven", mavenFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "username_Updated", mavenFeed.Username, "Feed username did not match expected value")
		assert.Equal(t, true, mavenFeed.Password.HasValue, "Feed password should be set")
		assert.Equal(t, "https://Updated.maven.apache.org/maven2/z", mavenFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, 7, mavenFeed.DownloadAttempts, "Feed download attempts did not match expected value")
		assert.Equal(t, 12, mavenFeed.DownloadRetryBackoffSeconds, "Feed download retry_backoff_seconds did not match expected value")

		return nil
	}
}
