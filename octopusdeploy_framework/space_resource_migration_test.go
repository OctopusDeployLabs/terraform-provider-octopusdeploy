package octopusdeploy_framework

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestSpaceResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testSpaceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: config,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testSpaceUpdated(t),
				),
			},
		},
	})
}

const config = `resource "octopusdeploy_space" "space_migration" {
						name                                 = "Test Space"
						space_managers_teams = ["teams-managers", "teams-administrators"]
					   }`

const updatedConfig = `resource "octopusdeploy_space" "space_migration" {
						  name = "Updated Test Space"
						  space_managers_teams = ["teams-managers"]
						  is_task_queue_stopped = true
					   }`

func testSpaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_space" {
			continue
		}

		space, err := octoClient.Spaces.GetByID(rs.Primary.ID)
		if err == nil && space != nil {
			return fmt.Errorf("space (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSpaceUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		spaceId := s.RootModule().Resources["octopusdeploy_space"+".space_migration"].Primary.ID
		space, err := octoClient.Spaces.GetByID(spaceId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve space by ID: %s", err)
		}

		assert.Equal(t, "Spaces-2", space.GetID(), "Space ID did not match expected value")
		assert.Equal(t, "Updated Test Space", space.Name, "Space name did not match expected value")
		assert.Equal(t, true, space.TaskQueueStopped, "Task Queue did not match expected value")
		assert.Equal(t, false, space.IsDefault, "IsDefault did not match expected value")
		assert.Contains(t, space.SpaceManagersTeams, "teams-managers", "Teams manager ID did not match expected value")

		return nil
	}
}
