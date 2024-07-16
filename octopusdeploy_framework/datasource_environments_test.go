package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceEnvironments(t *testing.T) {
	localName := acctest.RandStringFromCharSet(50, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_environments.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentsDataSourceID(prefix),
				),
				Config: testAccDataSourceEnvironmentsConfig(localName, take),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentsDataSourceID(prefix),
				),
				Config: testAccDataSourceEnvironmentsEmpty(localName),
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

func testAccDataSourceEnvironmentsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_environments %s" {}`, localName)
}
