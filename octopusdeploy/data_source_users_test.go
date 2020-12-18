package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceUsers(t *testing.T) {
	t.Parallel()

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_users.%s", localName)
	username := "d"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUsersConfig(localName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUsersDataSourceID(name),
					resource.TestCheckResourceAttrSet(name, "users.#"),
				)},
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
