package octopusdeploy_framework

import (
	"fmt"
	"os"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	internaltest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestRunbookResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	internaltest.SkipCI(t, "'octopusdeploy_runbook.runbook1' - expected NoOp, got action(s): [update]")
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		CheckDestroy: testRunbookDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.22.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: runbookConfig(name),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   runbookConfig(name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("octopusdeploy_runbook.runbook1", "NoOp"),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedRunbookConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testRunbookUpdated(t, name),
				),
			},
		},
	})
}

func runbookConfig(name string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "project1" {
		name = "project %[1]s"
		lifecycle_id = "Lifecycles-1"
		project_group_id = "ProjectGroups-1"
	}
	resource "octopusdeploy_runbook" "runbook1" {
		project_id = octopusdeploy_project.project1.id
		name = "runbook %[1]s"
	}
	`, name)
}

func updatedRunbookConfig(name string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "project1" {
		name = "project %[1]s"
		lifecycle_id = "Lifecycles-1"
		project_group_id = "ProjectGroups-1"
	}
	resource "octopusdeploy_runbook" "runbook1" {
		project_id = octopusdeploy_project.project1.id
		name = "runbook %[1]s"
		description = "description %[1]s"
		connectivity_policy {
			allow_deployments_to_no_targets = true
			exclude_unhealthy_targets = true
		}
		retention_policy {
			quantity_to_keep = 10
		}
	}
	`, name)
}

func testRunbookDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_runbook" {
			runbook, err := runbooks.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
			if err == nil && runbook != nil {
				return fmt.Errorf("runbook (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testRunbookUpdated(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		projectId := s.RootModule().Resources["octopusdeploy_project.project1"].Primary.ID
		runbookId := s.RootModule().Resources["octopusdeploy_runbook.runbook1"].Primary.ID
		runbook, err := runbooks.GetByID(octoClient, octoClient.GetSpaceID(), runbookId)
		if err != nil {
			return fmt.Errorf("failed to retrieve runbook by ID: %s", err)
		}

		assert.NotEmpty(t, runbook.ID, "Runbook ID did not match expected value")
		assert.Equal(t, runbook.ProjectID, projectId)
		assert.Equal(t, runbook.Description, fmt.Sprintf("description %s", name))
		assert.True(t, runbook.ConnectivityPolicy.AllowDeploymentsToNoTargets, "allow_deployments_to_no_targets should be true")
		assert.True(t, runbook.ConnectivityPolicy.ExcludeUnhealthyTargets, "exclude_unhealthy_targets should be true")
		assert.Equal(t, runbook.RunRetentionPolicy.QuantityToKeep, 10)

		return nil
	}
}
