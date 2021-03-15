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
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", value),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "14"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(localName, name, description, isSensitive, value, variableType),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", newValue),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "14"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(localName, name, description, isSensitive, newValue, variableType),
			},
			// {
			// 	ResourceName:      prefix,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testVariableBasic(localName string, name string, description string, isSensitive bool, value string, variableType string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "test-environment" {
	      name = "Test Environment (OK to Delete)"
		}

		resource "octopusdeploy_lifecycle" "test-lifecycle" {
		  name = "Test Lifecycle (OK to Delete)"
		}

		resource "octopusdeploy_project_group" "test-project-group" {
		  name = "Test Project Group (OK to Delete)"
		}

		resource "octopusdeploy_project" "test-project" {
		  lifecycle_id = octopusdeploy_lifecycle.test-lifecycle.id
		  name = "Test Project (OK to Delete)"
		  project_group_id = octopusdeploy_project_group.test-project-group.id
		}

		resource "octopusdeploy_channel" "test-channel" {
		  name = "Test Channel (OK to Delete)"
		  project_id = octopusdeploy_project.test-project.id
		}

		resource "octopusdeploy_variable" "%s" {
		  description = "%s"
		  is_sensitive = "%v"
		  name = "%s"
		  project_id = octopusdeploy_project.test-project.id
		  type = "%s"
		  value = "%s"

		  scope {
			channels = [octopusdeploy_channel.test-channel.id]
		    environments = [octopusdeploy_environment.test-environment.id]
		  }
		}`, localName, description, isSensitive, name, variableType, value)
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
