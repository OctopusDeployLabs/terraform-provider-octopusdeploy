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

func TestAzureSubscriptionAccountResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAzureSubscriptionAccountDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.43.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: azureSubscriptionAccountConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   azureSubscriptionAccountConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedAzureSubscriptionAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAzureSubscriptionAccountUpdated(t),
				),
			},
		},
	})
}

const azureSubscriptionAccountConfig = `
resource "octopusdeploy_azure_subscription_account" "azure_account" {
  name            = "Azure Subscription Account"
  subscription_id = "00000000-0000-0000-0000-000000000000"
 }`

const updatedAzureSubscriptionAccountConfig = ` 
resource "octopusdeploy_azure_subscription_account" "azure_account" {
  name            = "Update Azure Subscription Account"
  subscription_id = "00000000-0000-0000-0000-000000000000"
}`

func testAzureSubscriptionAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_account" {
			account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
			if err == nil && account != nil {
				return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAzureSubscriptionAccountUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountId := s.RootModule().Resources["octopusdeploy_aws_account.aws_account"].Primary.ID
		account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountId)
		if err != nil {
			return fmt.Errorf("failed to retrieve account by ID: %s", err)
		}

		azureAccount := account.(*accounts.AzureSubscriptionAccount)

		assert.NotEmpty(t, azureAccount.GetID(), "Account ID did not match expected value")
		assert.Equal(t, azureAccount.Name, "Updated test account", "Account name did not match expected value")

		return nil
	}
}
