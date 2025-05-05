package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
)

func TestAccOctopusDeployProcessTemplatedStep(t *testing.T) {
	paramDefaultValue := acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)
	scenario := newProcessTemplatedStepTestDependenciesConfiguration("template", paramDefaultValue)
	step := fmt.Sprintf("template_%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	defaultParameters := map[string]string{
		"Moves.One": paramDefaultValue,
	}
	requiredParameters := map[string]string{
		"Moves.Two": acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
	}
	allParameters := map[string]string{
		"Moves.One": acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
		"Moves.Two": acctest.RandStringFromCharSet(3, acctest.CharSetAlpha),
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProcessTemplatedStepConfiguration(scenario, step, requiredParameters),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessTemplatedStepAttributes(step, requiredParameters, defaultParameters),
					testCheckResourceProcessStepOfTypeExists("octopusdeploy_process_templated_step"),
				),
			},
			{
				Config: testAccProcessTemplatedStepConfiguration(scenario, step, allParameters),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceProcessTemplatedStepAttributes(step, allParameters, make(map[string]string)),
					testCheckResourceProcessStepOfTypeExists("octopusdeploy_process_templated_step"),
				),
			},
		},
	})
}

func testAccProcessTemplatedStepConfiguration(dependencies processTemplatedStepTestDependenciesConfiguration, step string, parameters map[string]string) string {
	var configurations []string
	for key, value := range parameters {
		configurations = append(configurations, fmt.Sprintf(`"%s" = "%s"`, key, value))
	}
	configuredParameters := strings.Join(configurations, "\n")

	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_process_templated_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "%s"
		  template_id = octopusdeploy_step_template.%s.id
		  template_version = octopusdeploy_step_template.%s.version
		  properties = {
			"Octopus.Action.TargetRoles" = "role-one"
		  }
		  parameters = {
			%s
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"
		  }
		}
		`,
		dependencies.config,
		step,
		dependencies.process,
		step,
		dependencies.template,
		dependencies.template,
		configuredParameters,
	)
}

func testCheckResourceProcessTemplatedStepAttributes(step string, parameters map[string]string, unmanagedParameters map[string]string) resource.TestCheckFunc {
	qualifiedName := fmt.Sprintf("octopusdeploy_process_templated_step.%s", step)

	assertions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(qualifiedName, "id"),
		resource.TestCheckResourceAttr(qualifiedName, "name", step),
		resource.TestCheckResourceAttr(qualifiedName, "type", "Octopus.Script"),
		resource.TestCheckResourceAttr(qualifiedName, "properties.Octopus.Action.TargetRoles", "role-one"),
		resource.TestCheckResourceAttr(qualifiedName, "template_properties.Octopus.Action.Script.ScriptSource", "Inline"),
		resource.TestCheckResourceAttr(qualifiedName, "template_properties.Octopus.Action.Script.Syntax", "Bash"),
		resource.TestCheckResourceAttr(qualifiedName, "template_properties.Octopus.Action.Script.ScriptBody", "echo '1.#{Moves.One} ... 2.#{Moves.Two} ... 3.?'"),
		resource.TestCheckResourceAttr(qualifiedName, "execution_properties.Octopus.Action.RunOnServer", "True"),
	}
	for key, value := range parameters {
		assertions = append(assertions, resource.TestCheckResourceAttr(qualifiedName, fmt.Sprintf("parameters.%s", key), value))
	}

	for key, value := range unmanagedParameters {
		assertions = append(assertions, resource.TestCheckResourceAttr(qualifiedName, fmt.Sprintf("unmanaged_parameters.%s", key), value))
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

type processTemplatedStepTestDependenciesConfiguration struct {
	process  string
	template string
	config   string
}

func newProcessTemplatedStepTestDependenciesConfiguration(scenario string, paramDefaultValue string) processTemplatedStepTestDependenciesConfiguration {
	projectGroup := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	project := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	process := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	template := fmt.Sprintf("%s_%s", scenario, acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	configuration := fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_step_template" "%s" {
		  action_type     = "Octopus.Script"
		  name            = "%s"
		  description     = "Template maintained by Terraform"
		  step_package_id = "Octopus.Script"
		  packages = []
		  parameters = [
			{
			  name      = "Moves.One"
			  id = "10001000-0000-0000-0000-100010001001"
			  default_value = "%s"
			  display_settings = {
				"Octopus.ControlType" : "SingleLineText"
			  }
			},
			{
			  name      = "Moves.Two"
			  id = "10001000-0000-0000-0000-100010001002"
			  display_settings = {
				"Octopus.ControlType" : "SingleLineText"
			  }
			},
		  ]	
		  properties = {
			"Octopus.Action.Script.ScriptBody" : "echo '1.#{Moves.One} ... 2.#{Moves.Two} ... 3.?'"
			"Octopus.Action.Script.ScriptSource" : "Inline"
			"Octopus.Action.Script.Syntax" : "Bash"
		  }
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
		template,
		template,
		paramDefaultValue,
		projectGroup,
		projectGroup,
		project,
		project,
		projectGroup,
		process,
		project,
	)

	return processTemplatedStepTestDependenciesConfiguration{
		process:  process,
		template: template,
		config:   configuration,
	}
}
