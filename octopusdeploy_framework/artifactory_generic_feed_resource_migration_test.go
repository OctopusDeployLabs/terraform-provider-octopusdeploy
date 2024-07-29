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

func (suite *IntegrationTestSuite) TestArtifactoryGenericFeedResource_UpgradeFromSDK_ToPluginFramework() {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy: testArtifactoryGenericFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: artifactoryGenericConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   artifactoryGenericConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   artifactoryGenericUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testArtifactoryGenericFeedUpdated(suite.T()),
				),
			},
		},
	})
}

const artifactoryGenericConfig = `resource "octopusdeploy_artifactory_generic_feed" "feed_artifactory_generic_migration" {
						  name                                 = "Helm"
						  feed_uri                             = "https://example.jfrog.io"
						  username                             = "username"
						  password                             = "password"
						  repository                     	   = "repo"
  						  layout_regex                   	   = "this is regex"
					   }`

const artifactoryGenericUpdatedConfig = `resource "octopusdeploy_artifactory_generic_feed" "feed_artifactory_generic_migration" {
						  name                                 = "Updated_Artifactory_Generic"
						  feed_uri                             = "https://example.jfrog.io/Updated"
						  username                             = "username_Updated"
						  password                             = "password_Updated"
						  repository                     	   = "repo_updated"
  						  layout_regex                   	   = "this is regex_updated"
					   }`

func testArtifactoryGenericFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_artifactory_generic_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testArtifactoryGenericFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_artifactory_generic_feed"+".feed_artifactory_generic_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		artifactoryGenericFeed := feed.(*feeds.ArtifactoryGenericFeed)

		assert.Regexp(t, "^Feeds\\-\\d+$", artifactoryGenericFeed.ID, "Feed ID did not match expected value")
		assert.Equal(t, "Updated_Artifactory_Generic", artifactoryGenericFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "username_Updated", artifactoryGenericFeed.Username, "Feed username did not match expected value")
		assert.Equal(t, true, artifactoryGenericFeed.Password.HasValue, "Feed password should be set")
		assert.Equal(t, "https://example.jfrog.io/Updated", artifactoryGenericFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, "repo_updated", artifactoryGenericFeed.Repository, "Feed repository should be set")
		assert.Equal(t, "this is regex_updated", artifactoryGenericFeed.LayoutRegex, "Feed layout regex did not match expected value")

		return nil
	}
}
