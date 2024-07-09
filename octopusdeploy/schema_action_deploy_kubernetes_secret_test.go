package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDeployKubernetesSecretAction(t *testing.T) {
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
				Config: testAccDeployKubernetesSecretAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployKubernetesSecretAction(),
				),
			},
		},
	})
}

func testAccDeployKubernetesSecretAction() string {
	return testAccBuildTestAction(`
		deploy_kubernetes_secret_action {
			name          = "Run Script"
			run_on_server = true
			secret_name   = "secret name"
			kubernetes_object_status_check_enabled = false

			secret_values = {
				"key-123" = "value-123",
				"key-321" = "value-321"
			}
    }
	`)
}

func testAccCheckDeployKubernetesSecretAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		process, err := getDeploymentProcess(s, octoClient)
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

		if action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"].Value != "False" {
			return fmt.Errorf("Kubernetes Object Status Check is incorrect: %s", action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"].Value)
		}

		return nil
	}
}
