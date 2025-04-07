package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strings"
	"testing"
)

func TestAccOctopusDeployProcessReplace(t *testing.T) {
	scenario := newProcessTestDependenciesConfiguration("replace")
	processName := fmt.Sprintf("replace_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentProcessConfiguration(scenario.config, processName, scenario.project1),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceDeploymentProcessAttributes(processName, scenario.project1),
					testCheckResourceDeploymentProcessBelongsToTheProject(processName, scenario.project1),
				),
			},
			{
				Config: testAccDeploymentProcessConfiguration(scenario.config, processName, scenario.project2),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceDeploymentProcessAttributes(processName, scenario.project2),
					testCheckResourceDeploymentProcessBelongsToTheProject(processName, scenario.project2),
				),
			},
			{
				Config: testAccRunbookProcessConfiguration(scenario.config, processName, scenario.project1, scenario.runbook),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceRunbookProcessAttributes(processName, scenario.project1, scenario.runbook),
					testCheckResourceRunbookProcessBelongsToTheRunbook(processName, scenario.runbook),
				),
			},
		},
	})
}

func testAccDeploymentProcessConfiguration(dependencies string, process string, project string) string {
	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		}
		`,
		dependencies,
		process,
		project,
	)
}

func testCheckResourceDeploymentProcessAttributes(processName string, projectName string) resource.TestCheckFunc {
	project := fmt.Sprintf("octopusdeploy_project.%s", projectName)
	process := fmt.Sprintf("octopusdeploy_process.%s", processName)

	assertions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(process, "id"),
		resource.TestCheckResourceAttrPair(process, "project_id", project, "id"),
		resource.TestCheckNoResourceAttr(process, "runbook_id"),
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

func testCheckResourceDeploymentProcessBelongsToTheProject(processName string, projectName string) resource.TestCheckFunc {
	project := fmt.Sprintf("octopusdeploy_project.%s", projectName)
	process := fmt.Sprintf("octopusdeploy_process.%s", processName)

	return func(state *terraform.State) error {
		processStateResource, ok := state.RootModule().Resources[process]
		if !ok {
			return fmt.Errorf("unable to find process resource: %s", process)
		}

		projectStateResource, ok := state.RootModule().Resources[project]
		if !ok {
			return fmt.Errorf("unable to find project resource: %s", project)
		}

		processId := processStateResource.Primary.ID
		projectId := projectStateResource.Primary.ID

		if strings.HasSuffix(processId, projectId) {
			return nil
		}

		return fmt.Errorf("expected process id '%s' to belong to project id '%s'", processId, projectId)
	}
}

func testAccRunbookProcessConfiguration(dependencies string, process string, project string, runbook string) string {
	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		  runbook_id  = octopusdeploy_runbook.%s.id
		}
		`,
		dependencies,
		process,
		project,
		runbook,
	)
}

func testCheckResourceRunbookProcessAttributes(processName string, projectName string, runbookName string) resource.TestCheckFunc {
	project := fmt.Sprintf("octopusdeploy_project.%s", projectName)
	runbook := fmt.Sprintf("octopusdeploy_runbook.%s", runbookName)
	process := fmt.Sprintf("octopusdeploy_process.%s", processName)

	assertions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(process, "id"),
		resource.TestCheckResourceAttrPair(process, "project_id", project, "id"),
		resource.TestCheckResourceAttrPair(process, "runbook_id", runbook, "id"),
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

func testCheckResourceRunbookProcessBelongsToTheRunbook(processName string, runbookName string) resource.TestCheckFunc {
	runbook := fmt.Sprintf("octopusdeploy_runbook.%s", runbookName)
	process := fmt.Sprintf("octopusdeploy_process.%s", processName)

	return func(state *terraform.State) error {
		processStateResource, ok := state.RootModule().Resources[process]
		if !ok {
			return fmt.Errorf("unable to find process resource: %s", process)
		}

		runbookStateResource, ok := state.RootModule().Resources[runbook]
		if !ok {
			return fmt.Errorf("unable to find runbook resource: %s", runbook)
		}

		processId := processStateResource.Primary.ID
		runbookId := runbookStateResource.Primary.ID

		if strings.HasSuffix(processId, runbookId) {
			return nil
		}

		return fmt.Errorf("expected process id '%s' to belong to runbook id '%s'", processId, runbookId)
	}
}

type processTestDependenciesConfiguration struct {
	projectGroup string
	project1     string
	project2     string
	runbook      string
	config       string
}

func newProcessTestDependenciesConfiguration(scenario string) processTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project1 := fmt.Sprintf("%s_1_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project2 := fmt.Sprintf("%s_2_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	runbook := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
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

		resource "octopusdeploy_runbook" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		  name = "%s"
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
		runbook,
		project1,
		runbook,
	)

	return processTestDependenciesConfiguration{
		projectGroup: projectGroup,
		project1:     project1,
		project2:     project2,
		runbook:      runbook,
		config:       configuration,
	}
}
