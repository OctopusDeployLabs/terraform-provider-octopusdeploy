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

func TestAwsElasticContainerRegistryResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")
	os.Setenv("ECR_ACCESS_KEY", "")
	os.Setenv("ECR_SECRET_KEY", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAwsElasticContainerRegistryFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: awsElasticContainerRegistryConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   awsElasticContainerRegistryConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   awsElasticContainerRegistryUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAwsElasticContainerRegistryFeedUpdated(t),
				),
			},
		},
	})
}

const awsElasticContainerRegistryConfig = `resource "octopusdeploy_aws_elastic_container_registry" "feed_aws_elastic_container_registry_migration" {
						  name                                 = "AWS ECR feed"
						  region                     	   	   = "us-east-1"
  						  secret_key                   	       = "secret"
						  access_key                   	       = "access key"
					   }`

const awsElasticContainerRegistryUpdatedConfig = `resource "octopusdeploy_aws_elastic_container_registry" "feed_aws_elastic_container_registry_migration" {
						  name                                 = "AWS ECR feed Updated"
						  region                     	   	   = "us-east-2"
  						  secret_key                   	       = "secret updated"
						  access_key                   	       = "access key updated"
					   }`

func testAwsElasticContainerRegistryFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_elastic_container_registry" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAwsElasticContainerRegistryFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_aws_elastic_container_registry"+".feed_aws_elastic_container_registry_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		awsElasticContainerRegistry := feed.(*feeds.AwsElasticContainerRegistry)

		assert.Equal(t, "Feeds-1001", awsElasticContainerRegistry.ID, "Feed ID did not match expected value")
		assert.Equal(t, "AWS ECR feed Updated", awsElasticContainerRegistry.Name, "Feed name did not match expected value")
		assert.Equal(t, "access key updated", awsElasticContainerRegistry.AccessKey, "Feed access key should be set")
		assert.Equal(t, "us-east-2", awsElasticContainerRegistry.Region, "Feed Region did not match expected value")
		assert.Equal(t, true, awsElasticContainerRegistry.SecretKey.HasValue, "Feed secret key should be set")

		return nil
	}
}
