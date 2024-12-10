package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
	"time"
)

func TestAccDataSourceDeploymentFreezes(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	freezeName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	dataSourceName := "data.octopusdeploy_deployment_freezes.test_freeze"

	startTime := time.Now().AddDate(1, 0, 0).UTC()
	endTime := startTime.Add(24 * time.Hour)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName, startTime, endTime, false, false, projectName, environmentName, tenantName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "partial_name", freezeName),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.name", freezeName),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.project_environment_scope.%", "0"),
				),
			},
			{
				Config: testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName, startTime, endTime, true, false, projectName, environmentName, tenantName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.project_environment_scope.%", "1"),
				),
			},
			{
				Config: testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName, startTime, endTime, true, true, projectName, environmentName, tenantName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.0.tenant_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "deployment_freezes.0.tenant_project_environment_scope.0.environment_id"),
					resource.TestCheckResourceAttr(dataSourceName, "deployment_freezes.0.project_environment_scope.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceDeploymentFreezesConfig(spaceName, freezeName string, startTime, endTime time.Time, includeProject bool, includeTenant bool, projectName, environmentName, tenantName string) string {
	baseConfig := fmt.Sprintf(`
resource "octopusdeploy_space" "test_space" {
    name                  = "%s"
    is_default            = false
    is_task_queue_stopped = false
    description           = "Test space for deployment freeze datasource"
    space_managers_teams  = ["teams-administrators"]
}

resource "octopusdeploy_environment" "test_environment" {
    name        = "%s"
    description = "Test environment"
    space_id    = octopusdeploy_space.test_space.id
}

resource "octopusdeploy_project_group" "test_project_group" {
    name        = "Test Project Group"
    description = "Test project group for deployment freeze datasource"
    space_id    = octopusdeploy_space.test_space.id
}

resource "octopusdeploy_lifecycle" "test_lifecycle" {
    name     = "Test Lifecycle"
    space_id = octopusdeploy_space.test_space.id
}

resource "octopusdeploy_project" "test_project" {
    name           = "%s"
    lifecycle_id   = octopusdeploy_lifecycle.test_lifecycle.id
    project_group_id = octopusdeploy_project_group.test_project_group.id
    space_id       = octopusdeploy_space.test_space.id
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
`, spaceName, environmentName, projectName, freezeName, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	if includeProject {
		projectConfig := fmt.Sprintf(`
resource "octopusdeploy_deployment_freeze_project" "test_project_scope" {
    deploymentfreeze_id = octopusdeploy_deployment_freeze.test_freeze.id
    project_id = octopusdeploy_project.test_project.id
    environment_ids = [octopusdeploy_environment.test_environment.id]
}
`)
		baseConfig = baseConfig + projectConfig
	}

	if includeTenant {
		tenantConfig := fmt.Sprintf(`
resource "octopusdeploy_tenant" "test_tenant" {
    name     = "%s"
    space_id = octopusdeploy_space.test_space.id
}

resource "octopusdeploy_tenant_project" "test_tenant_project" {
    tenant_id = octopusdeploy_tenant.test_tenant.id
    project_id = octopusdeploy_project.test_project.id
    environment_ids = [octopusdeploy_environment.test_environment.id]
    space_id = octopusdeploy_space.test_space.id
}

resource "octopusdeploy_deployment_freeze_tenant" "test_tenant_scope" {
    deploymentfreeze_id = octopusdeploy_deployment_freeze.test_freeze.id
    tenant_id = octopusdeploy_tenant.test_tenant.id
    project_id = octopusdeploy_project.test_project.id
    environment_id = octopusdeploy_environment.test_environment.id

    depends_on = [
        octopusdeploy_tenant_project.test_tenant_project
    ]
}
`, tenantName)
		baseConfig = baseConfig + tenantConfig
	}

	datasourceConfig := `
data "octopusdeploy_deployment_freezes" "test_freeze" {
    partial_name = "%s"
    skip         = 0
    take         = 1
    depends_on   = [`

	deps := []string{"octopusdeploy_deployment_freeze.test_freeze"}
	if includeProject {
		deps = append(deps, "octopusdeploy_deployment_freeze_project.test_project_scope")
	}
	if includeTenant {
		deps = append(deps, "octopusdeploy_deployment_freeze_tenant.test_tenant_scope")
	}

	datasourceConfig += "\n        " + strings.Join(deps, ",\n        ") + "\n    "

	datasourceConfig += `]
}

output "octopus_space_id" {
    value = octopusdeploy_space.test_space.id
}

output "octopus_freeze_id" {
    value = data.octopusdeploy_deployment_freezes.test_freeze.deployment_freezes[0].id
}
`
	var config = baseConfig + fmt.Sprintf(datasourceConfig, freezeName)
	return config
}
