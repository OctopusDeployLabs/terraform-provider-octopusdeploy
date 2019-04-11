package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	const tfVarPrefix = "octopusdeploy_variable.foovar"
	const tfVarName = "tf-var-1"
	const tfVarDesc = "Terraform testing module variable"
	const tfVarValue = "abcd-123456"
	const projectName = "Funky Monkey Var Test"
	const lifeCycleID = "Lifecycles-1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVariableBasic(projectName, lifeCycleID, tfVarName, tfVarDesc, tfVarValue),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployVariableExists(tfVarPrefix),
					resource.TestCheckResourceAttr(
						tfVarPrefix, "name", tfVarName),
					resource.TestCheckResourceAttr(
						tfVarPrefix, "description", tfVarDesc),
					resource.TestCheckResourceAttr(
						tfVarPrefix, "value", tfVarValue),
				),
			},
		},
	})
}

func testVariableBasic(projectName, projectLifecycleID, name, description, value string) string {
	config := fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			name           = "%s"
			lifecycle_id    = "%s"
			project_group_id = "${octopusdeploy_project_group.foo.id}"
		}

		resource "octopusdeploy_variable" "foovar" {
			project_id  = "${octopusdeploy_project.foo.id}"
			name        = "%s"
			description = "%s"
			type        = "String"
			value       = "%s"
		}
		`,
		projectName, projectLifecycleID, name, description, value,
	)
	fmt.Println(config)
	return config
}

func testOctopusDeployVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsVarHelper(s, client)
	}
}

func existsVarHelper(s *terraform.State, client *octopusdeploy.Client) error {
	projID := s.RootModule().Resources["octopusdeploy_project.foo"].Primary.ID
	varID := s.RootModule().Resources["octopusdeploy_variable.foovar"].Primary.ID

	if _, err := client.Variable.GetByID(projID, varID); err != nil {
		return fmt.Errorf("Received an error retrieving variable %s", err)
	}

	return nil
}

func testOctopusDeployVariableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyVarHelper(s, client)
}

func destroyVarHelper(s *terraform.State, client *octopusdeploy.Client) error {
	projID := s.RootModule().Resources["octopusdeploy_project.foo"].Primary.ID
	varID := s.RootModule().Resources["octopusdeploy_variable.foovar"].Primary.ID

	if _, err := client.Variable.DeleteSingle(projID, varID); err != nil {
		if err == octopusdeploy.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving variable %s", err)
	}
	return fmt.Errorf("Variable still exists")
}
