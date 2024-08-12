package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceLifecycles(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := "Default Lifecycle"
	resourceName := "data.octopusdeploy_lifecycles.lifecycle_default_lifecycle"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLifecyclesConfig(spaceName, lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "partial_name", lifecycleName),
					resource.TestCheckResourceAttr(resourceName, "lifecycles.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "lifecycles.0.id"),
					resource.TestCheckResourceAttr(resourceName, "lifecycles.0.name", lifecycleName),
					testAccCheckOutputExists("octopus_space_id"),
					testAccCheckOutputExists("octopus_lifecycle_id"),
				),
			},
		},
	})
}

func testAccDataSourceLifecyclesConfig(spaceName, lifecycleName string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_space" "octopus_project_space_test" {
  name                  = "%s"
  is_default            = false
  is_task_queue_stopped = false
  description           = "Test space for lifecycles datasource"
  space_managers_teams  = ["teams-administrators"]
}

data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "%s"
  space_id     = octopusdeploy_space.octopus_project_space_test.id
  skip         = 0
  take         = 1
  depends_on   = [octopusdeploy_space.octopus_project_space_test]
}

output "octopus_space_id" {
  value = octopusdeploy_space.octopus_project_space_test.id
}

output "octopus_lifecycle_id" {
  value = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
}
`, spaceName, lifecycleName)
}
