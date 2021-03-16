package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployRunScriptAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRunScriptAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRunScriptAction(),
				),
			},
		},
	})
}

func testAccRunScriptAction() string {
	return testAccBuildTestAction(`
		run_script_action {
            name = "Run Script"
            run_on_server = true
			script_file_name = "Test.ps1"
			script_parameters = "-Test 1"
			script_source = "Package"
			variable_substitution_in_files = "test.json"

			primary_package {
				package_id = "MyPackage"
				feed_id = "feeds-builtin"
			}

			package {
				acquisition_location = "Server"
				feed_id = octopusdeploy_nuget_feed.testing.id
				name = "package2"
				package_id = "package2"
				extract_during_deployment = false
				properties = {
				  "Extract" = "false"
				}
			}
        }
	`)
}

func testAccCheckRunScriptAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.Script" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Script.ScriptFileName"] != "Test.ps1" {
			return fmt.Errorf("ScriptFileName is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptFileName"])
		}

		if action.Properties["Octopus.Action.Script.ScriptParameters"] != "-Test 1" {
			return fmt.Errorf("ScriptSource is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptParameters"])
		}

		if action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] != "test.json" {
			return fmt.Errorf("TargetFiles is incorrect: %s", action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"])
		}

		return nil
	}
}
