package octopusdeploy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccEnvironment(localName string, name string, description string, allowDynamicInfrastructure bool, sortOrder int, useGuidedFailure bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		allow_dynamic_infrastructure = "%v"
		description                  = "%s"
		name                         = "%s"
		sort_order                   = %v
		use_guided_failure           = "%v"
	}`, localName, allowDynamicInfrastructure, description, name, sortOrder, useGuidedFailure)
}

func testAccEnvironmentCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_environment" {
			continue
		}

		if environment, err := octoClient.Environments.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("environment (%s) still exists", environment.GetID())
		}
	}

	return nil
}
