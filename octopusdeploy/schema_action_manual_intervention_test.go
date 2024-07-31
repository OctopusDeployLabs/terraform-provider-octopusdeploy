package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployManualInterventionAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
			sort_order = 1
			instructions = "Approve Me"
			responsible_teams = "A Team"
		}
	`)
}

func testAccCheckManualInterventionAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		process, err := getDeploymentProcess(s, octoClient)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.Manual" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Manual.Instructions"].Value != "Approve Me" {
			return fmt.Errorf("Instructions is incorrect: %s", action.Properties["Octopus.Action.Manual.Instructions"].Value)
		}

		if action.Properties["Octopus.Action.Manual.ResponsibleTeamIds"].Value != "A Team" {
			return fmt.Errorf("ResponsibleTeamIds is incorrect: %s", action.Properties["Octopus.Action.Manual.ResponsibleTeamIds"].Value)
		}

		return nil
	}
}
