package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccOctopusDeployProcessChildStepManualIntervention(t *testing.T) {
	scenario := newProcessChildStepTestDependenciesConfiguration("child")
	step := fmt.Sprintf("child_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	instruction1 := fmt.Sprintf(`Create: %s`, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))
	instruction2 := fmt.Sprintf(`Update: %s`, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessChildStepManualConfiguration(scenario, step, instruction1),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessChildStepManualAttributes(scenario, step, instruction1),
					testCheckResourceProcessChildStepExists(),
				),
			},
			{
				Config: testAccProcessChildStepManualConfiguration(scenario, step, instruction2),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessChildStepManualAttributes(scenario, step, instruction2),
					testCheckResourceProcessChildStepExists(),
				),
			},
		},
	})
}

func testAccProcessChildStepManualConfiguration(scenario processChildStepTestDependenciesConfiguration, step string, instructions string) string {
	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_process_child_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  parent_id = octopusdeploy_process_step.%s.id
		  name = "%s"
		  type = "Octopus.Manual"
		  execution_properties = {
    		"Octopus.Action.RunOnServer" = "True"
    		"Octopus.Action.Manual.Instructions" = "%s"
    		"Octopus.Action.Manual.BlockConcurrentDeployments" = "True"
    		"Octopus.Action.Manual.ResponsibleTeamIds" = "teams-managers"
		  }
		}
		`,
		scenario.config,
		step,
		scenario.process,
		scenario.parent,
		step,
		instructions,
	)
}

func testCheckResourceProcessChildStepManualAttributes(scenario processChildStepTestDependenciesConfiguration, step string, instructions string) resource.TestCheckFunc {
	process := fmt.Sprintf("octopusdeploy_process.%s", scenario.process)
	parent := fmt.Sprintf("octopusdeploy_process_step.%s", scenario.parent)
	childStep := fmt.Sprintf("octopusdeploy_process_child_step.%s", step)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(childStep, "id"),
		resource.TestCheckResourceAttrPair(childStep, "process_id", process, "id"),
		resource.TestCheckResourceAttrPair(childStep, "parent_id", parent, "id"),
		resource.TestCheckResourceAttr(childStep, "name", step),
		resource.TestCheckResourceAttr(childStep, "type", "Octopus.Manual"),
		resource.TestCheckResourceAttr(childStep, "execution_properties.Octopus.Action.RunOnServer", "True"),
		resource.TestCheckResourceAttr(childStep, "execution_properties.Octopus.Action.Manual.Instructions", instructions),
		resource.TestCheckResourceAttr(childStep, "execution_properties.Octopus.Action.Manual.BlockConcurrentDeployments", "True"),
		resource.TestCheckResourceAttr(childStep, "execution_properties.Octopus.Action.Manual.ResponsibleTeamIds", "teams-managers"),
	)
}

type processChildStepTestDependenciesConfiguration struct {
	process string
	parent  string
	config  string
}

func newProcessChildStepTestDependenciesConfiguration(scenario string) processChildStepTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	process := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	parentStep := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	configuration := fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}
		resource "octopusdeploy_project_group" "%s" {
		  name        = "%s"
		  description = "Test process step"
		}

		resource "octopusdeploy_project" "%s" {
		  name                                 = "%s"
		  description                          = "Test process step"
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
			"Octopus.Action.RunOnServer" = "True"
			"Octopus.Action.Script.ScriptSource" = "Inline"
			"Octopus.Action.Script.Syntax"       = "PowerShell"
			"Octopus.Action.Script.ScriptBody" = "Write-Host 'step with children'"
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
		parentStep,
		process,
		parentStep,
	)

	return processChildStepTestDependenciesConfiguration{
		process: process,
		parent:  parentStep,
		config:  configuration,
	}
}

func testCheckResourceProcessChildStepExists() resource.TestCheckFunc {
	return testCheckResourceProcessChildStepOfTypeExists("octopusdeploy_process_child_step")
}

func testCheckResourceProcessChildStepOfTypeExists(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type == resourceType {
				actionId := r.Primary.ID
				stepId := r.Primary.Attributes["parent_id"]
				processId := r.Primary.Attributes["process_id"]
				process, err := deployments.GetDeploymentProcessByID(octoClient, octoClient.GetSpaceID(), processId)
				if err != nil {
					return fmt.Errorf("expected process with id '%s' to exist: %s", processId, err)
				}

				step, stepExists := deploymentProcessWrapper{process}.FindStepByID(stepId)
				if !stepExists {
					return fmt.Errorf("expected process (%s) to contain step (%s)", processId, stepId)
				}

				_, actionExists := findActionFromProcessStepByID(step, actionId)
				if !actionExists {
					return fmt.Errorf("expected process (%s) to contain step (%s) with action (%s)", processId, stepId, actionId)
				}
			}
		}
		return nil
	}
}
