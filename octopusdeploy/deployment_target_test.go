package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testDeploymentTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		deploymentTargetID := s.RootModule().Resources[resourceName].Primary.ID
		if _, err := client.Machines.GetByID(deploymentTargetID); err != nil {
			return fmt.Errorf("error retrieving deployment target: %s", err)
		}

		return nil
	}
}

func testDeploymentTargetCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_deployment_target" {
			continue
		}

		_, err := client.Machines.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("deployment target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
