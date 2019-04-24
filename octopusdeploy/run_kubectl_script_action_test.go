package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployRunKubectlScriptAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRunKubectlScriptAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRunKubectlScriptAction(),
				),
			},
		},
	})
}

func testAccRunKubectlScriptAction() string {
	return testAccBuildTestAction(`
		run_kubectl_script_action {
            name = "Run Script"
            run_on_server = true
			
			primary_package {
				package_id = "MyPackage"
				feed_id = "feeds-builtin"
			}
			
			script_file_name = "Test.ps1"
			script_parameters = "-Test 1"
        }
	`)
}

func testAccCheckRunKubectlScriptAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.KubernetesRunScript" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Script.ScriptFileName"] != "Test.ps1" {
			return fmt.Errorf("ScriptFileName is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptFileName"])
		}

		if action.Properties["Octopus.Action.Script.ScriptParameters"] != "-Test 1" {
			return fmt.Errorf("ScriptSource is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptParameters"])
		}

		return nil
	}
}
