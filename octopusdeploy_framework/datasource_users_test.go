package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccDataSourceUsers(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_users.%s", localName)
	username := "d"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUsersDataSourceID(name),
					resource.TestCheckResourceAttrSet(name, "users.#"),
				),
				Config: testAccDataSourceUsersConfig(localName, username),
			},
		},
	})
}

func testAccCheckUsersDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Users data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Users source ID not set")
		}
		return nil
	}
}

func testAccDataSourceUsersConfig(localName string, username string) string {
	return fmt.Sprintf(`data "octopusdeploy_users" "%s" {
		filter = "%s"
	}`, localName, username)
}
