package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployRunKubectlScriptAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

			namespace = "test-namespace"
    }
	`)
}

func testAccCheckRunKubectlScriptAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.KubernetesRunScript" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.KubernetesContainers.Namespace"].Value != "test-namespace" {
			return fmt.Errorf("Kubernetes namespace is incorrect: %s", action.Properties["Octopus.Action.KubernetesContainers.Namespace"].Value)
		}

		if action.Properties["Octopus.Action.Script.ScriptFileName"].Value != "Test.ps1" {
			return fmt.Errorf("ScriptFileName is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptFileName"].Value)
		}

		if action.Properties["Octopus.Action.Script.ScriptParameters"].Value != "-Test 1" {
			return fmt.Errorf("ScriptSource is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptParameters"].Value)
		}

		return nil
	}
}
