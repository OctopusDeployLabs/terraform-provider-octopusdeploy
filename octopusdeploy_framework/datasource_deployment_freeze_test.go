package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
	"time"
)

func TestAccDataSourceDeploymentFreezes(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	freezeName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "data.octopusdeploy_deployment_freezes.test_freeze"

	// Use future dates for the freeze period, ensuring UTC timezone
	startTime := time.Now().AddDate(1, 0, 0).UTC()
	endTime := startTime.Add(24 * time.Hour)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "partial_name", freezeName),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_freezes.0.id"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.name", freezeName),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.start", startTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.end", endTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_freezes.0.project_environment_scope.%"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.recurring_schedule.type", "DaysPerWeek"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.recurring_schedule.unit", "24"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.recurring_schedule.end_type", "AfterOccurrences"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.recurring_schedule.end_after_occurrences", "5"),
					resource.TestCheckResourceAttr(resourceName, "deployment_freezes.0.recurring_schedule.days_of_week.#", "3"),
					testAccCheckOutputExists("octopus_space_id"),
					testAccCheckOutputExists("octopus_freeze_id"),
				),
			},
		},
	})
}

func testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName string, startTime, endTime time.Time) string {
	return fmt.Sprintf(`
resource "octopusdeploy_space" "test_space" {
	name                  = "%s"
	is_default            = false
	is_task_queue_stopped = false
	description           = "Test space for deployment freeze datasource"
	space_managers_teams  = ["teams-administrators"]
}

resource "octopusdeploy_deployment_freeze" "test_freeze" {
	name  = "%s"
	start = "%s"
	end   = "%s"
	
	recurring_schedule = {
		type = "DaysPerWeek"
		unit = 24
		end_type = "AfterOccurrences"
		end_after_occurrences = 5
		days_of_week = ["Monday", "Wednesday", "Friday"]
	}

	depends_on = [octopusdeploy_space.test_space]
}

data "octopusdeploy_deployment_freezes" "test_freeze" {
	ids          = null
	partial_name = "%s"
	skip         = 0
	take         = 1
	depends_on   = [octopusdeploy_deployment_freeze.test_freeze]
}

output "octopus_space_id" {
	value = octopusdeploy_space.test_space.id
}

output "octopus_freeze_id" {
	value = data.octopusdeploy_deployment_freezes.test_freeze.deployment_freezes[0].id
}
`, spaceName, freezeName, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), freezeName)
}
