package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
	return testAccBuildTestAction(`
		manual_intervention_action {
			name = "Test"
			instructions = "Approve Me"
			responsible_teams = "A Team"
		}
	`)
}

func testAccCheckManualInterventionAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
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

		return nil
	}
}
