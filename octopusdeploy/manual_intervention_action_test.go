package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployManualInterventionAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManualInterventionAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManualInterventionAction(),
				),
			},
		},
	})
}

func testAccManualInterventionAction() string {
	return `
		resource "octopusdeploy_lifecycle" "test" {
			name = "Test Lifecycle"
		}

		resource "octopusdeploy_project_group" "test" {
			name = "Test Group"
		}

		resource "octopusdeploy_project" "test" {
			name             = "Test Project"
			lifecycle_id     = "${octopusdeploy_lifecycle.test.id}"
			project_group_id = "${octopusdeploy_project_group.test.id}"
		}

		resource "octopusdeploy_deployment_process" "test" {
			project_id = "${octopusdeploy_project.test.id}"

			step {
				name = "Test"

				manual_intervention_action {
					name = "Test"
					instructions = "Approve Me"
					responsible_teams = "A Team"
				}
			}
		}
		`
}

func testAccCheckManualInterventionAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client);
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.Manual" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Manual.Instructions"] != "Approve Me" {
			return fmt.Errorf("Instructions is incorrect: %s", action.Properties["Octopus.Action.Manual.Instructions"])
		}

		if action.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] != "A Team" {
			return fmt.Errorf("ResponsibleTeamIds is incorrect: %s", action.Properties["Octopus.Action.Manual.ResponsibleTeamIds"])
		}


		return nil;
	}
}
