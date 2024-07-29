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

func (suite *IntegrationTestSuite) TestNugetFeed_UpgradeFromSDK_ToPluginFramework() {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testNugetFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: nugetConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   nugetConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedNugetConfig,
				Check: resource.ComposeTestCheckFunc(
					testNugetFeedUpdated(t),
				),
			},
		},
	})
}

const nugetConfig = `resource "octopusdeploy_nuget_feed" "feed_nuget_migration" {
						  name                                 = "Nuget"
						  feed_uri                             = "https://api.nuget.org/v3/index.json"
						  username                             = "username"
						  password                             = "password"
						  is_enhanced_mode					   = false
						  download_attempts                    = 6
						  download_retry_backoff_seconds       = 11
					   }`

const updatedNugetConfig = `resource "octopusdeploy_nuget_feed" "feed_nuget_migration" {
						  name                                 = "Updated Nuget"
						  feed_uri                             = "https://api.nuget.org/v4/index.json"
						  username                             = "username_Updated"
						  password                             = "password_Updated"
 	 					  is_enhanced_mode					   = true
						  download_attempts                    = 7
						  download_retry_backoff_seconds       = 12
					   }`

func testNugetFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_nuget_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testNugetFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_nuget_feed"+".feed_nuget_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		nugetFeed := feed.(*feeds.NuGetFeed)

		assert.Equal(t, "Feeds-1001", nugetFeed.ID, "Feed ID did not match expected value")
		assert.Equal(t, "Updated Nuget", nugetFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "username_Updated", nugetFeed.Username, "Feed username did not match expected value")
		assert.Equal(t, true, nugetFeed.Password.HasValue, "Feed password should be set")
		assert.Equal(t, true, nugetFeed.EnhancedMode, "Feed enhanced mode should be set")
		assert.Equal(t, "https://api.nuget.org/v4/index.json", nugetFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, 7, nugetFeed.DownloadAttempts, "Feed download attempts did not match expected value")
		assert.Equal(t, 12, nugetFeed.DownloadRetryBackoffSeconds, "Feed download retry_backoff_seconds did not match expected value")

		return nil
	}
}
