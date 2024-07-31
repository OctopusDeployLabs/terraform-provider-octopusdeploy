package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func (suite *IntegrationTestSuite) TestEnvironmentResource_UpgradeFromSDK_ToPluginFramework() {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		CheckDestroy: testEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.22.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: environmentConfig(name),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   environmentConfig(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updateEnvironmentConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testEnvironment(suite.T(), name),
				),
			},
		},
	})
}

func environmentConfig(name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "environment1" {
		name = "%s"
		sort_order = 1
	}`, name)
}

func updateEnvironmentConfig(name string) string {
	return fmt.Sprintf(
		`resource "octopusdeploy_environment" "environment1" {
			name = "%s"
			description = "%s"
			sort_order = 1
		}`, name, name)
}

func testEnvironmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_environment" {
			continue
		}

		environment, err := octoClient.Environments.GetByID(rs.Primary.ID)
		if err == nil && environment != nil {
			return fmt.Errorf("environment (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testEnvironment(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		environmentId := s.RootModule().Resources["octopusdeploy_environment.environment1"].Primary.ID
		environment, err := octoClient.Environments.GetByID(environmentId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve environment by ID: %s", err)
		}

		assert.NotEmpty(t, environment.ID, "Environment ID did not match expected value")
		assert.Equal(t, name, environment.Name, "Environment name did not match expected value")
		assert.Equal(t, name, environment.Description, "Environment description did not match expected value")

		return nil
	}
}
