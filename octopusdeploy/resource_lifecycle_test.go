package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccLifecycle(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
		release_retention_policy {
			unit             = "Days"
			quantity_to_keep = 30
		}

		tentacle_retention_policy {
			unit             = "Days"
			quantity_to_keep = 30
		}
	}`, localName, name)
}

func testAccCheckLifecycleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := existsHelperLifecycle(s, octoClient); err != nil {
			return err
		}
		return nil
	}
}

func existsHelperLifecycle(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_lifecycle" {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}
func testAccLifecycleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_lifecycle" {
			continue
		}

		lifecycle, err := octoClient.Lifecycles.GetByID(rs.Primary.ID)
		if err == nil && lifecycle != nil {
			return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
