package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceSpaces(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_spaces.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpacesConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpacesDataSourceID(name),
				)},
		},
	})
}

func testAccCheckSpacesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Spaces data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Spaces source ID not set")
		}
		return nil
	}
}

func testAccDataSourceSpacesConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_spaces" "%s" {
		take = %v
	}`, localName, take)
}
