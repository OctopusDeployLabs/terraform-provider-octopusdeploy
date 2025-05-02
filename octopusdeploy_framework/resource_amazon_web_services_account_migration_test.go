package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAmazonWebServicesAccountResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAmazonWebServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.43.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: amazonWebServicesAccountConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   amazonWebServicesAccountConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedAmazonWebServicesAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAmazonWebServicesAccountUpdated(t),
				),
			},
		},
	})
}

const amazonWebServicesAccountConfig = `
resource "octopusdeploy_environment" "environment1" {
	name = "env-1"
	description = "environment1"
}

resource "octopusdeploy_environment" "environment2" {
	name = "env-2"
	description = "environment2"
} 

resource "octopusdeploy_aws_account" "aws_account" {
  name         = "Test account" 
  access_key   = "access-key"
  secret_key   = "secret-key" 
  environments = [octopusdeploy_environment.environment1.id, octopusdeploy_environment.environment2.id] 
}`

const updatedAmazonWebServicesAccountConfig = ` 
resource "octopusdeploy_environment" "environment1" {
	name = "env-1"
	description = "environment1"
}

resource "octopusdeploy_environment" "environment2" {
	name = "env-2"
	description = "environment2"
} 

resource "octopusdeploy_aws_account" "aws_account" {
  name         = "Updated test account" 
  access_key   = "access-key"
  secret_key   = "secret-key" 
  environments = [octopusdeploy_environment.environment1.id]
}`

func testAmazonWebServicesAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_account" {
			account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
			if err == nil && account != nil {
				return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAmazonWebServicesAccountUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountId := s.RootModule().Resources["octopusdeploy_aws_account.aws_account"].Primary.ID
		account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountId)
		if err != nil {
			return fmt.Errorf("failed to retrieve account by ID: %s", err)
		}

		awsAccount := account.(*accounts.AmazonWebServicesAccount)

		assert.NotEmpty(t, awsAccount.GetID(), "Account ID did not match expected value")
		assert.Equal(t, awsAccount.Name, "Updated test account", "Account name did not match expected value")

		return nil
	}
}
