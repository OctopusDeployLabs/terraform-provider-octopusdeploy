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

func TestAccOctopusDeployProcessStepsOrder(t *testing.T) {
	scenario := newProcessStepsOrderTestDependenciesConfiguration("order")
	order := fmt.Sprintf("order_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	orderedSteps1 := []string{scenario.step1, scenario.step2, scenario.step3}
	orderedSteps2 := []string{scenario.step3, scenario.step1, scenario.step2}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessStepsOrderConfiguration(scenario, order, orderedSteps1),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessStepsOrderAttributes(scenario, order, orderedSteps1),
					testCheckResourceProcessStepsOrderExists(t, orderedSteps1),
				),
			},
			{
				Config: testAccProcessStepsOrderConfiguration(scenario, order, orderedSteps2),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessStepsOrderAttributes(scenario, order, orderedSteps2),
					testCheckResourceProcessStepsOrderExists(t, orderedSteps2),
				),
			},
		},
	})
}

func testAccProcessStepsOrderConfiguration(scenario processStepsOrderTestDependenciesConfiguration, orderResource string, steps []string) string {
	stepReferences := make([]string, len(steps))
	for i, step := range steps {
		stepReferences[i] = fmt.Sprintf("octopusdeploy_process_step.%s.id,", step)
	}
	orderedSteps := strings.Join(stepReferences, "\n")

	return fmt.Sprintf(`
		%s
		resource octopusdeploy_process_steps_order "%s" {
		  process_id = octopusdeploy_process.%s.id
		  steps = [
			%s
		  ]
		}
		`,
		scenario.config,
		orderResource,
		scenario.process,
		orderedSteps,
	)
}

func testCheckResourceProcessStepsOrderAttributes(scenario processStepsOrderTestDependenciesConfiguration, name string, steps []string) resource.TestCheckFunc {
	order := fmt.Sprintf("octopusdeploy_process_steps_order.%s", name)
	process := fmt.Sprintf("octopusdeploy_process.%s", scenario.process)

	assertions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(order, "id", order, "process_id"),
		resource.TestCheckResourceAttrPair(order, "id", process, "id"),
	}

	for i, expected := range steps {
		step := fmt.Sprintf("octopusdeploy_process_step.%s", expected)
		stepPosition := fmt.Sprintf("steps.%d", i)
		assertions = append(assertions, resource.TestCheckResourceAttrPair(order, stepPosition, step, "id"))
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

type processStepsOrderTestDependenciesConfiguration struct {
	process string
	step1   string
	step2   string
	step3   string
	config  string
}

func newProcessStepsOrderTestDependenciesConfiguration(scenario string) processStepsOrderTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	process := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	step1 := fmt.Sprintf("%s_1_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	step2 := fmt.Sprintf("%s_2_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	step3 := fmt.Sprintf("%s_3_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

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

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "%s"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-two"
		  }
		  type = "Octopus.Script"
		  execution_properties = {
			"Octopus.Action.Script.ScriptBody" = "."
		  }
		}

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "%s"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-three"
		  }
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
		step1,
		process,
		step1,
		step2,
		process,
		step2,
		step3,
		process,
		step3,
	)

	return processStepsOrderTestDependenciesConfiguration{
		process: process,
		step1:   step1,
		step2:   step2,
		step3:   step3,
		config:  configuration,
	}
}

func testCheckResourceProcessStepsOrderExists(t *testing.T, expectedSteps []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_process_steps_order" {
				processId := r.Primary.ID
				process, err := deployments.GetDeploymentProcessByID(octoClient, octoClient.GetSpaceID(), processId)
				if err != nil {
					return fmt.Errorf("expected process with id '%s' to exist: %s", processId, err)
				}

				actualOrder := make([]string, len(process.Steps))
				for i, step := range process.Steps {
					actualOrder[i] = step.Name
				}

				assert.Equal(t, expectedSteps, actualOrder, "Persisted process steps should be ordered as expected")
			}
		}
		return nil
	}
}
