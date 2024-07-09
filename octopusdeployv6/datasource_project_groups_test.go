package octopusdeployv6

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceProjectGroups(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_project_groups.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		// It would be nice to use the ProtoV6ProviderFactories approach but I've had issues with the init functionality
		// For now I'm sticking with the external providers and terraformrc over-write behaviour we use in the existing integration tests.
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectGroupsDataSourceID(prefix),
				),
				Config: testAccDataSourceProjectGroupsConfig(localName, take),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectGroupsDataSourceID(prefix),
				),
				Config: testAccDataSourceProjectGroupsEmpty(localName),
			},
		},
	})
}

func testAccCheckProjectGroupsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find ProjectGroups data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot ProjectGroups source ID not set")
		}
		return nil
	}
}

func testAccDataSourceProjectGroupsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_project_groups" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceProjectGroupsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_project_groups" "%s" {}`, localName)
}
