package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceSpaces(t *testing.T) {
	spaceID := "Spaces-1"
	resourceName := "data.octopusdeploy_spaces.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpacesConfig(spaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ids.0", spaceID),
					resource.TestCheckResourceAttr(resourceName, "skip", "0"),
					resource.TestCheckResourceAttr(resourceName, "take", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "spaces.0.id"),
					testAccCheckOutputExists("octopus_space_id"),
					resource.TestCheckOutput("octopus_space_id", spaceID),
				),
			},
		},
	})
}

func testAccDataSourceSpacesConfig(spaceID string) string {
	tfConfig := fmt.Sprintf(`
		data "octopusdeploy_spaces" "test" {
		  ids  = ["%s"]
		  skip = 0
		  take = 1
		}
		
		output "octopus_space_id" {
		  value = data.octopusdeploy_spaces.test.spaces[0].id
		}
		`, spaceID)
	return tfConfig
}

func testAccCheckOutputExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Outputs[name]
		if !ok {
			return fmt.Errorf("output %s not found", name)
		}
		return nil
	}
}
