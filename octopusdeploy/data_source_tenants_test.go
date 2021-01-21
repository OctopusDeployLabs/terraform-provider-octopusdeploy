package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTenants(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_tenants.%s", localName)
	skip := acctest.RandIntRange(0, 100)
	take := acctest.RandIntRange(0, 100)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantsDataSourceID(name),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttr(name, "skip", strconv.Itoa(skip)),
					resource.TestCheckResourceAttr(name, "take", strconv.Itoa(take)),
				),
				Config: testAccDataSourceTenantsConfig(localName, skip, take),
			},
		},
	})
}

func testAccCheckTenantsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find tenants data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot tenants source ID not set")
		}
		return nil
	}
}

func testAccDataSourceTenantsConfig(localName string, skip int, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_tenants" "%s" {
	  skip = %v
	  take = %v
	}`, localName, skip, take)
}
