package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccOctopusDeployProcessStepRunScript(t *testing.T) {
	scenario := newProcessStepTestDependenciesConfiguration("basic")
	step := fmt.Sprintf("basic_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	createScript := fmt.Sprintf(`Write-Host 'create: %s'`, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))
	updateScript := fmt.Sprintf(`Write-Host 'update: %s'`, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessStepRunScriptConfiguration(scenario.config, scenario.process, step, createScript),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessStepRunScriptAttributes(step, createScript),
					testCheckResourceProcessStepExists(),
				),
			},
			{
				Config: testAccProcessStepRunScriptConfiguration(scenario.config, scenario.process, step, updateScript),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessStepRunScriptAttributes(step, updateScript),
					testCheckResourceProcessStepExists(),
				),
			},
		},
	})
}

func testAccProcessStepRunScriptConfiguration(dependencies string, process string, step string, scriptBody string) string {
	return fmt.Sprintf(`
		%s
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
			"Octopus.Action.Script.ScriptBody" = "%s"
		  }
		}
		`,
		dependencies,
		step,
		process,
		step,
		scriptBody,
	)
}

func testCheckResourceProcessStepRunScriptAttributes(step string, script string) resource.TestCheckFunc {
	qualifiedName := fmt.Sprintf("octopusdeploy_process_step.%s", step)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(qualifiedName, "id"),
		resource.TestCheckResourceAttr(qualifiedName, "name", step),
		resource.TestCheckResourceAttr(qualifiedName, "type", "Octopus.Script"),
		resource.TestCheckResourceAttr(qualifiedName, "properties.Octopus.Action.TargetRoles", "role-one"),
		resource.TestCheckResourceAttr(qualifiedName, "execution_properties.Octopus.Action.RunOnServer", "True"),
		resource.TestCheckResourceAttr(qualifiedName, "execution_properties.Octopus.Action.Script.ScriptSource", "Inline"),
		resource.TestCheckResourceAttr(qualifiedName, "execution_properties.Octopus.Action.Script.Syntax", "PowerShell"),
		resource.TestCheckResourceAttr(qualifiedName, "execution_properties.Octopus.Action.Script.ScriptBody", script),
	)
}

type processStepTestDependenciesConfiguration struct {
	process string
	config  string
}

func newProcessStepTestDependenciesConfiguration(scenario string) processStepTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	process := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
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
		`,
		projectGroup,
		projectGroup,
		project,
		project,
		projectGroup,
		process,
		project,
	)

	return processStepTestDependenciesConfiguration{
		process: process,
		config:  configuration,
	}
}

func testCheckResourceProcessStepExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_process_step" {
				stepId := r.Primary.ID
				processId := r.Primary.Attributes["process_id"]
				process, processError := deployments.GetDeploymentProcessByID(octoClient, octoClient.GetSpaceID(), processId)
				if processError != nil {
					return fmt.Errorf("expected process with id '%s' to exist: %s", processId, processError)
				}

				_, stepExists := deploymentProcessWrapper{process}.FindStepByID(stepId)
				if !stepExists {
					return fmt.Errorf("expected process (%s) to contain step (%s)", processId, stepId)
				}
			}
		}
		return nil
	}
}
