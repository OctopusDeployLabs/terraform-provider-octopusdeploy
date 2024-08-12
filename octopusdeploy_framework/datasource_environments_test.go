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

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	//take := 10

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
			// 			{
			// 				Check: resource.ComposeTestCheckFunc(
			// 					testAccCheckEnvironmentsDataSourceID(prefix),
			// 					resource.TestCheckResourceAttr(prefix, "environments.#", "3"),
			// 				),
			// 				Config: fmt.Sprintf(`%s

			// %s`,
			// 					createTestAccDataSourceEnvironmentsConfig(environmentLocalName, environmentName),
			// 					testAccDataSourceEnvironmentsConfig(localName, take),
			// 				),
			// 			},
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
					createTestAccDataSourceEnvironmentsConfig(environmentLocalName, environmentName),
					testAccDataSourceEnvironmentByNameConfig(localName, environmentName),
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

func testAccDataSourceEnvironmentsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceEnvironmentByNameConfig(localName string, name string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {
		name = "%s"	
	}`, localName, name)
}

func testAccDataSourceEnvironmentsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments" "%s" {}`, localName)
}

func createTestAccDataSourceEnvironmentsConfig(localName string, name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%[1]s" {
			name = "%[2]s"
		}
		
		resource "octopusdeploy_environment" "%[1]s-1" {
			name = "%[2]s-1"
		}

		resource "octopusdeploy_environment" "%[1]s-2" {
			name = "%[2]s-2"
		}
	`, localName, name)
}
