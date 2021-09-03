package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDeployKuberentesSecretAction(t *testing.T) {
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
				Config: testAccDeployKuberentesSecretAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployKuberentesSecretAction(),
				),
			},
		},
	})
}

func testAccDeployKuberentesSecretAction() string {
	return testAccBuildTestAction(`
		deploy_kubernetes_secret_action {
			name          = "Run Script"
			run_on_server = true
			secret_name   = "secret name"

			secret_values = {
				"key-123" = "value-123",
				"key-321" = "value-321"
			}
    }
	`)
}

func testAccCheckDeployKuberentesSecretAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.KubernetesDeploySecret" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.KubernetesContainers.SecretName"].Value != "secret name" {
			return fmt.Errorf("SecretName is incorrect: %s", action.Properties["Octopus.Action.KubernetesContainers.SecretName"].Value)
		}

		if action.Properties["Octopus.Action.KubernetesContainers.SecretValues"].Value != `{"key-123":"value-123","key-321":"value-321"}` {
			return fmt.Errorf("SecretValue is incorrect: %s", action.Properties["Octopus.Action.KubernetesContainers.SecretValues"].Value)
		}

		return nil
	}
}
