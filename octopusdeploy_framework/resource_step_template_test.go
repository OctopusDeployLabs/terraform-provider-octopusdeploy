package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

type stepTemplatePackagePropsTestData struct {
	extract       string
	purpose       string
	selectionMode string
}

type stepTemplatePackageTestData struct {
	packageID          string
	acquisitonLocation string
	feedID             string
	name               string
	properties         stepTemplatePackagePropsTestData
}

type stepTemplateParamTestData struct {
	defaultValue    string
	displaySettings map[string]string
	helpText        string
	label           string
	name            string
	id              string
}

type stepTemplateTestData struct {
	localName     string
	prefix        string
	actionType    string
	name          string
	description   string
	stepPackageID string
	packages      []stepTemplatePackageTestData
	parameters    []stepTemplateParamTestData
	properties    map[string]string
}

func TestAccOctopusStepTemplateBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_step_template." + localName
	data := stepTemplateTestData{
		localName:     localName,
		prefix:        prefix,
		actionType:    "Octopus.Script",
		name:          acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
		description:   acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		stepPackageID: "Octopus.Script",
		packages: []stepTemplatePackageTestData{
			{
				packageID:          "force",
				acquisitonLocation: "Server",
				feedID:             "feeds-builtin",
				name:               "mypackage",
				properties: stepTemplatePackagePropsTestData{
					extract:       "True",
					purpose:       "",
					selectionMode: "immediate",
				},
			},
		},
		parameters: []stepTemplateParamTestData{
			{
				defaultValue: "Hello World",
				displaySettings: map[string]string{
					"Octopus.ControlType": "SingleLineText",
				},
				helpText: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				label:    acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				name:     acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				id:       "621e1584-cdf3-4b67-9204-fc82430c908c",
			},
			{
				defaultValue: "Hello Earth",
				displaySettings: map[string]string{
					"Octopus.ControlType": "SingleLineText",
				},
				helpText: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				label:    acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				name:     acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				id:       "cd731d21-669a-42e1-81af-048681fd5c69",
			},
		},
		properties: map[string]string{
			"Octopus.Action.Script.ScriptBody":   "echo 'Hello World'",
			"Octopus.Action.Script.ScriptSource": "Inline",
			"Octopus.Action.Script.Syntax":       "Bash",
		},
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testStepTemplateDestroy(s, localName) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testStepTemplateRunScriptBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.name),
				),
			},
			{
				Config: testStepTemplateRunScriptUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.name+"-updated"),
				),
			},
		},
	})
}

func testStepTemplateRunScriptBasic(data stepTemplateTestData) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_step_template" "%s" {
			action_type     = "%s"
			name            = "%s"
			description     = "%s"
		  	step_package_id = "%s"
		  	packages        = [
				{
					package_id = "%s"
					acquisition_location = "%s"
					feed_id = "%s"
					name = "%s"
					properties = {
						extract = "%s"
						purpose = "%s"
						selection_mode = "%s"
				  	}
				}
			]
		  	parameters = [
				{
			  		default_value = "%s"
			  		display_settings = {
						"Octopus.ControlType" : "%s"
			  		}
			  		help_text = "%s"
			  		label     = "%s"
			  		name      = "%s"
			  		id = "%s"
				},
				{
			  		default_value = "%s"
			  		display_settings = {
						"Octopus.ControlType" : "%s"
			  		}
			  		help_text = "%s"
			  		label     = "%s"
			  		name      = "%s"
			  		id = "%s"
				},
		  	]
		  	properties = {
				"Octopus.Action.Script.ScriptBody" : "%s"
				"Octopus.Action.Script.ScriptSource" : "%s"
				"Octopus.Action.Script.Syntax" : "%s"
		  	}
		}
`,
		data.localName,
		data.actionType,
		data.name,
		data.description,
		data.stepPackageID,
		data.packages[0].packageID,
		data.packages[0].acquisitonLocation,
		data.packages[0].feedID,
		data.packages[0].name,
		data.packages[0].properties.extract,
		data.packages[0].properties.purpose,
		data.packages[0].properties.selectionMode,
		data.parameters[0].defaultValue,
		data.parameters[0].displaySettings["Octopus.ControlType"],
		data.parameters[0].helpText,
		data.parameters[0].label,
		data.parameters[0].name,
		data.parameters[0].id,
		data.parameters[1].defaultValue,
		data.parameters[1].displaySettings["Octopus.ControlType"],
		data.parameters[1].helpText,
		data.parameters[1].label,
		data.parameters[1].name,
		data.parameters[1].id,
		data.properties["Octopus.Action.Script.ScriptBody"],
		data.properties["Octopus.Action.Script.ScriptSource"],
		data.properties["Octopus.Action.Script.Syntax"],
	)
}

func testStepTemplateRunScriptUpdate(data stepTemplateTestData) string {
	data.name = data.name + "-updated"

	return testStepTemplateRunScriptBasic(data)
}

func testStepTemplateDestroy(s *terraform.State, localName string) error {
	var actionTemplateID string

	if rs, ok := s.RootModule().Resources[localName]; ok {
		if rs.Type != "octopusdeploy_step_template" {
			return fmt.Errorf("resource has unexpected type: %s", rs.Type)
		}
		actionTemplateID = rs.Primary.ID
	} else {
		return fmt.Errorf("resource octopusdeploy_template.%s not found in state", localName)
	}

	actionTemplate, err := actiontemplates.GetByID(octoClient, octoClient.GetSpaceID(), actionTemplateID)
	if err == nil {
		if actionTemplate != nil {
			return fmt.Errorf("step template (%s) still exists", actionTemplate.Name)
		}
	}

	return nil
}
