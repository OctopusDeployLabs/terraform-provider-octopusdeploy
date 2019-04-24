package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployApplyTerraformAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplyTerraformAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplyTerraformAction(),
				),
			},
		},
	})
}

func testAccApplyTerraformAction() string {
	return testAccBuildTestAction(`
		apply_terraform_action {
            name = "Apply Terraform"
            run_on_server = true
			
			primary_package {
				package_id = "MyPackage"
				feed_id = "feeds-builtin"
			}
			
			additional_init_params = "Init params"
        }
	`)
}

func testAccCheckApplyTerraformAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.TerraformApply" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Terraform.AdditionalInitParams"] != "Init params" {
			return fmt.Errorf("AdditionalInitParams is incorrect: %s", action.Properties["Octopus.Action.Terraform.AdditionalInitParams"])
		}

		return nil
	}
}
