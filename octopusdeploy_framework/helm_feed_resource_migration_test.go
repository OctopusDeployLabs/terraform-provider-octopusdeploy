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

func TestHelmFeedResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testHelmFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: helmConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   helmConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   helmUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testHelmFeedUpdated(t),
				),
			},
		},
	})
}

const helmConfig = `resource "octopusdeploy_helm_feed" "feed_helm_migration" {
						  name                                 = "Helm"
						  feed_uri                             = "https://charts.helm.sh/stable/"
						  username                             = "username"
						  password                             = "password"
					   }`

const helmUpdatedConfig = `resource "octopusdeploy_helm_feed" "feed_helm_migration" {
						  name                                 = "Updated_Helm"
						  feed_uri                             = "https://charts.helm.sh/stableUpdated/"
						  username                             = "username_Updated"
						  password                             = "password_Updated"
					   }`

func testHelmFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_helm_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testHelmFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_helm_feed"+".feed_helm_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		helmFeed := feed.(*feeds.HelmFeed)

		assert.Equal(t, "Feeds-1001", helmFeed.ID, "Feed ID did not match expected value")
		assert.Equal(t, "Updated_Helm", helmFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "username_Updated", helmFeed.Username, "Feed username did not match expected value")
		assert.Equal(t, true, helmFeed.Password.HasValue, "Feed password should be set")
		assert.Equal(t, "https://charts.helm.sh/stableUpdated/", helmFeed.FeedURI, "Feed URI did not match expected value")

		return nil
	}
}
