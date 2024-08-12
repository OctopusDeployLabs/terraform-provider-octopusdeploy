package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceEnvironments(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_environments.%s", localName)

	spaceName := "env_datasource_test" // acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	take := 10

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentsDataSourceID(prefix),
				),
				Config: testAccDataSourceEnvironmentsEmpty(localName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentsDataSourceID(prefix),
					resource.TestCheckResourceAttr(prefix, "environments.#", "3"),
				),
				Config: fmt.Sprintf(`%s
				
				%s`,
					createTestAccDataSourceEnvironmentsConfig(spaceName, environmentLocalName, environmentName),
					testAccDataSourceEnvironmentsConfig(localName, take, spaceName, environmentLocalName),
				),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentsDataSourceID(prefix),
					resource.TestCheckResourceAttr(prefix, "name", environmentName),
					resource.TestCheckResourceAttr(prefix, "environments.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "environments.0.id"),
					resource.TestCheckResourceAttr(prefix, "environments.0.name", environmentName),
				),
				Config: fmt.Sprintf(`%s

			%s`,
					createTestAccDataSourceEnvironmentsConfig(spaceName, environmentLocalName, environmentName),
					testAccDataSourceEnvironmentByNameConfig(localName, environmentName, spaceName, environmentLocalName),
				),
			},
		},
	})
}

func testAccCheckEnvironmentsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Environments data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Environments source ID not set")
		}
		return nil
	}
}

func createTestAccDataSourceEnvironmentsConfig(spaceName string, localName string, name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "%[1]s" {
			name                  = "%[1]s"
			is_default            = false
			is_task_queue_stopped = false
			description           = "Test space for environments datasource"
			space_managers_teams  = ["teams-administrators"]
		}

		resource "octopusdeploy_environment" "%[2]s" {
			name = "%[3]s"
			space_id = octopusdeploy_space.%[1]s.id
		}
		
		resource "octopusdeploy_environment" "%[2]s-1" {
			name = "%[3]s-1"
			space_id = octopusdeploy_space.%[1]s.id
		}

		resource "octopusdeploy_environment" "%[2]s-2" {
			name = "%[3]s-2"
			space_id = octopusdeploy_space.%[1]s.id
		}
	`, spaceName, localName, name)
}

func testAccDataSourceEnvironmentsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {}`, localName)
}

func testAccDataSourceEnvironmentsConfig(localName string, take int, spaceName string, environmentLocalName string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {
		take = %v
		space_id = octopusdeploy_space.%s.id
		depends_on = [octopusdeploy_environment.%s, octopusdeploy_environment.%[4]s-1, octopusdeploy_environment.%[4]s-2]
	}`, localName, take, spaceName, environmentLocalName)
}

func testAccDataSourceEnvironmentByNameConfig(localName string, name string, spaceName string, environmentLocalName string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {
		name = "%s"	
		space_id = octopusdeploy_space.%s.id
		depends_on = [octopusdeploy_environment.%s, octopusdeploy_environment.%[4]s-1, octopusdeploy_environment.%[4]s-2]
	}`, localName, name, spaceName, environmentLocalName)
}
