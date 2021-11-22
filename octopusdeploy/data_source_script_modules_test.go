package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceScriptModules(t *testing.T) {
	t.Parallel()

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_script_modules.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScriptModulesConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScriptModulesDataSourceID(name),
				)},
		},
	})
}

func testAccCheckScriptModulesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find script modules data source: %s", n)
		}

		if len(rs.Primary.ID) <= 0 {
			return fmt.Errorf("script modules source ID not set")
		}
		return nil
	}
}

func testAccDataSourceScriptModulesConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_script_modules" "%s" {
		take = %v
	}`, localName, take)
}
