package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_variable." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isSensitive := false
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	variableType := "String"

	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testVariableDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "owner_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", value),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "2"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, value, variableType),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "owner_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", newValue),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "2"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, newValue, variableType),
			},
			// {
			// 	ResourceName:      prefix,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func TestAccOctopusDeployVariableSchemaValidation(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_variable." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testVariableDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccVariableSchemaValidation(localName, name, projectID),
			},
		},
	})
}

func testVariableBasic(environmentLocalName string,
	environmentName string,
	lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, channelLocalName string, channelName string, localName string, name string, description string, isSensitive bool, value string, variableType string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
	      name = "%s"
		}

		resource "octopusdeploy_lifecycle" "%s" {
		  name = "%s"
		}

		resource "octopusdeploy_project_group" "%s" {
		  name = "%s"
		}

		resource "octopusdeploy_project" "%s" {
		  lifecycle_id     = octopusdeploy_lifecycle.%s.id
		  name             = "%s"
		  project_group_id = octopusdeploy_project_group.%s.id
		}

		resource "octopusdeploy_channel" "%s" {
		  name       = "%s"
		  project_id = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_variable" "%s" {
		  description  = "%s"
		  is_sensitive = "%v"
		  name         = "%s"
		  owner_id     = octopusdeploy_project.%s.id
		  type         = "%s"
		  value        = "%s"

		  scope {
			channels     = [octopusdeploy_channel.%s.id]
		    environments = [octopusdeploy_environment.%s.id]
		  }
		}`,
		environmentLocalName,
		environmentName,
		lifecycleLocalName,
		lifecycleName,
		projectGroupLocalName,
		projectGroupName,
		projectLocalName,
		lifecycleLocalName,
		projectName,
		projectGroupLocalName,
		channelLocalName,
		channelName,
		projectLocalName,
		localName,
		description,
		isSensitive,
		name,
		projectLocalName,
		variableType,
		value,
		channelLocalName,
		environmentLocalName,
	)
}

func testAccVariableSchemaValidation(localName string, name string, projectID string) string {
	return fmt.Sprintf(`resource "octopusdeploy_variable" "%s" {
		name       = "%s"
		project_id = "%s"
		type       = "String"
		value      = "1"
	  }`, localName, name, projectID)
}

func testOctopusDeployVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var projectID string
		var variableID string

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_project" {
				projectID = r.Primary.ID
			}

			if r.Type == "octopusdeploy_variable" {
				projectID = r.Primary.Attributes["project_id"]
				variableID = r.Primary.ID
			}
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if _, err := client.Variables.GetByID(projectID, variableID); err != nil {
			return fmt.Errorf("error retrieving variable %s", err)
		}

		return nil
	}
}

func testVariableDestroy(s *terraform.State) error {
	var projectID string
	var variableID string

	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_project" {
			projectID = r.Primary.ID
		}

		if r.Type == "octopusdeploy_variable" {
			variableID = r.Primary.ID
		}
	}

	client := testAccProvider.Meta().(*octopusdeploy.Client)
	variable, err := client.Variables.GetByID(projectID, variableID)
	if err == nil {
		if variable != nil {
			return fmt.Errorf("variable (%s) still exists", variableID)
		}
	}

	return nil
}
