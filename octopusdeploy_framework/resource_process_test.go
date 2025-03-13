package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strings"
	"testing"
)

func TestAccOctopusDeployProcessBasicReplace(t *testing.T) {
	scenario := newProcessTestDependenciesConfiguration("basic")
	processName := fmt.Sprintf("basic_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	processQualifiedName := fmt.Sprintf("octopusdeploy_process.%s", processName)

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessBasicConfiguration(scenario.config, processName, scenario.project1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(processQualifiedName, "id"),
					testCheckResourceProcessBelongsToTheProject(processQualifiedName, scenario.project1),
				),
			},
			{
				Config: testAccProcessBasicConfiguration(scenario.config, processName, scenario.project2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(processQualifiedName, "id"),
					testCheckResourceProcessBelongsToTheProject(processQualifiedName, scenario.project2),
				),
			},
		},
	})
}

func testAccProcessBasicConfiguration(dependencies string, processName string, ownerName string) string {
	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_process" "%s" {
		  owner_id  = %s.id
		}
		`,
		dependencies,
		processName,
		ownerName,
	)
}

func testCheckResourceProcessBelongsToTheProject(process string, project string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		processResource, ok := state.RootModule().Resources[process]
		if !ok {
			return fmt.Errorf("unable to find process resource: %s", process)
		}

		projectResource, ok := state.RootModule().Resources[project]
		if !ok {
			return fmt.Errorf("unable to find project resource: %s", process)
		}

		processId := processResource.Primary.ID
		projectId := projectResource.Primary.ID

		if strings.HasSuffix(processId, projectId) {
			return nil
		}

		return fmt.Errorf("expected process id '%s' to belong to project id '%s'", processId, projectId)
	}
}

type processTestDependenciesConfiguration struct {
	projectGroup string
	project1     string
	project2     string
	config       string
}

func newProcessTestDependenciesConfiguration(scenario string) processTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project1 := fmt.Sprintf("%s_1_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project2 := fmt.Sprintf("%s_2_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	configuration := fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}
		resource "octopusdeploy_project_group" "%s" {
		  name        = "%s"
		  description = "Test of basic process"
		}

		resource "octopusdeploy_project" "%s" {
		  name                                 = "%s"
		  description                          = "Test of basic process"
		  default_guided_failure_mode          = "EnvironmentDefault"
		  tenanted_deployment_participation    = "Untenanted"
		  project_group_id                     = octopusdeploy_project_group.%s.id
		  lifecycle_id                         = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  included_library_variable_sets       = []
		}

		resource "octopusdeploy_project" "%s" {
		  name                                 = "%s"
		  description                          = "Test of basic process"
		  default_guided_failure_mode          = "EnvironmentDefault"
		  tenanted_deployment_participation    = "Untenanted"
		  project_group_id                     = octopusdeploy_project_group.%s.id
		  lifecycle_id                         = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  included_library_variable_sets       = []
		}
		`,
		projectGroup,
		projectGroup,
		project1,
		project1,
		projectGroup,
		project2,
		project2,
		projectGroup,
	)

	return processTestDependenciesConfiguration{
		projectGroup: fmt.Sprintf("%s.%s", "octopusdeploy_project_group", projectGroup),
		project1:     fmt.Sprintf("%s.%s", "octopusdeploy_project", project1),
		project2:     fmt.Sprintf("%s.%s", "octopusdeploy_project", project2),
		config:       configuration,
	}
}
