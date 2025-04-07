package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAccOctopusDeployProcessChildStepsOrder(t *testing.T) {
	scenario := newProcessChildStepsOrderTestDependenciesConfiguration("children")
	order := fmt.Sprintf("children_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	orderedSteps1 := []string{scenario.child1, scenario.child2, scenario.child3}
	orderedSteps2 := []string{scenario.child3, scenario.child1, scenario.child2}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessChildStepsOrderConfiguration(scenario, order, orderedSteps1),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessChildStepsOrderAttributes(scenario, order, orderedSteps1),
					testCheckResourceProcessChildStepsOrderExists(t, scenario.parent, orderedSteps1),
				),
			},
			{
				Config: testAccProcessChildStepsOrderConfiguration(scenario, order, orderedSteps2),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessChildStepsOrderAttributes(scenario, order, orderedSteps2),
					testCheckResourceProcessChildStepsOrderExists(t, scenario.parent, orderedSteps2),
				),
			},
		},
	})
}

func testAccProcessChildStepsOrderConfiguration(scenario processChildStepsOrderTestDependenciesConfiguration, orderResource string, children []string) string {
	references := make([]string, len(children))
	for i, step := range children {
		references[i] = fmt.Sprintf("octopusdeploy_process_child_step.%s.id,", step)
	}
	orderedChildren := strings.Join(references, "\n")

	return fmt.Sprintf(`
		%s
		resource octopusdeploy_process_child_steps_order "%s" {
		  process_id = octopusdeploy_process.%s.id
		  parent_id = octopusdeploy_process_step.%s.id
		  children = [
			%s
		  ]
		}
		`,
		scenario.config,
		orderResource,
		scenario.process,
		scenario.parent,
		orderedChildren,
	)
}

func testCheckResourceProcessChildStepsOrderAttributes(scenario processChildStepsOrderTestDependenciesConfiguration, name string, children []string) resource.TestCheckFunc {
	order := fmt.Sprintf("octopusdeploy_process_child_steps_order.%s", name)
	process := fmt.Sprintf("octopusdeploy_process.%s", scenario.process)
	parent := fmt.Sprintf("octopusdeploy_process_step.%s", scenario.parent)

	assertions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(order, "id", order, "parent_id"),
		resource.TestCheckResourceAttrPair(order, "process_id", process, "id"),
		resource.TestCheckResourceAttrPair(order, "parent_id", parent, "id"),
	}

	for i, expected := range children {
		step := fmt.Sprintf("octopusdeploy_process_child_step.%s", expected)
		stepPosition := fmt.Sprintf("children.%d", i)
		assertions = append(assertions, resource.TestCheckResourceAttrPair(order, stepPosition, step, "id"))
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

type processChildStepsOrderTestDependenciesConfiguration struct {
	process string
	parent  string
	child1  string
	child2  string
	child3  string
	config  string
}

func newProcessChildStepsOrderTestDependenciesConfiguration(scenario string) processChildStepsOrderTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	process := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	parent := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	child1 := fmt.Sprintf("%s_1_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	child2 := fmt.Sprintf("%s_2_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	child3 := fmt.Sprintf("%s_3_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	configuration := fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}
		resource "octopusdeploy_project_group" "%s" {
		  name        = "%s"
		  description = "Test process steps ordering"
		}

		resource "octopusdeploy_project" "%s" {
		  name                                 = "%s"
		  description                          = "Test process steps ordering"
		  default_guided_failure_mode          = "EnvironmentDefault"
		  tenanted_deployment_participation    = "Untenanted"
		  project_group_id                     = octopusdeploy_project_group.%s.id
		  lifecycle_id                         = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  included_library_variable_sets       = []
		}

		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "%s"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-one"
		  }
		  type = "Octopus.Script"
		  execution_properties = {
			"Octopus.Action.Script.ScriptBody" = "."
		  }
		}

		resource "octopusdeploy_process_child_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  parent_id   = octopusdeploy_process_step.%s.id
		  name = "%s"
		  type = "Octopus.Script"
		  execution_properties = {
			"Octopus.Action.Script.ScriptBody" = "."
		  }
		}

		resource "octopusdeploy_process_child_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  parent_id   = octopusdeploy_process_step.%s.id
		  name = "%s"
		  type = "Octopus.Script"
		  execution_properties = {
			"Octopus.Action.Script.ScriptBody" = "."
		  }
		}

		resource "octopusdeploy_process_child_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  parent_id   = octopusdeploy_process_step.%s.id
		  name = "%s"
		  type = "Octopus.Script"
		  execution_properties = {
			"Octopus.Action.Script.ScriptBody" = "."
		  }
		}
		`,
		projectGroup,
		projectGroup,
		project,
		project,
		projectGroup,
		process,
		project,
		parent,
		process,
		parent,
		child1,
		process,
		parent,
		child1,
		child2,
		process,
		parent,
		child2,
		child3,
		process,
		parent,
		child3,
	)

	return processChildStepsOrderTestDependenciesConfiguration{
		process: process,
		parent:  parent,
		child1:  child1,
		child2:  child2,
		child3:  child3,
		config:  configuration,
	}
}

func testCheckResourceProcessChildStepsOrderExists(t *testing.T, parent string, children []string) resource.TestCheckFunc {
	// Based on assumption that first action (embedded within step in configuration) have same name as a parent step
	expectedActions := append([]string{parent}, children...)

	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_process_child_steps_order" {
				stepId := r.Primary.ID
				processId := r.Primary.Attributes["process_id"]
				process, processError := deployments.GetDeploymentProcessByID(octoClient, octoClient.GetSpaceID(), processId)
				if processError != nil {
					return fmt.Errorf("expected process with id '%s' to exist: %s", processId, processError)
				}

				step, stepExists := deploymentProcessWrapper{process}.FindStepByID(stepId)
				if !stepExists {
					return fmt.Errorf("expected process (%s) to contain step (%s)", processId, stepId)
				}

				actualOrder := make([]string, len(step.Actions))
				for i, action := range step.Actions {
					actualOrder[i] = action.Name
				}

				assert.Equal(t, expectedActions, actualOrder, "Persisted process step actions should be ordered as expected")
			}
		}
		return nil
	}
}
