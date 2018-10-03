package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	const tfVarPrefix = "octopusdeploy_variable.foo"
	const tfVarProjectID = "Projects-1"
	const tfVarName = "tf-var-1"
	const tfVarDesc = "Terraform testing module variable"
	const tfVarValue = "abcd-123456"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVariableBasic(tfVarProjectID, tfVarName, tfVarDesc, tfVarValue),
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

func testVariableBasic(projectID, name, description, value string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_variable" "foo" {
			project_id  = "%s"
			name        = "%s"
			description = "%s"
			type        = "String"
			value       = "%s"
		}
		`,
		projectID, name, description, value,
	)
}

func testOctopusDeployVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsVarHelper(s, client)
	}
}

func existsVarHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Variable.GetByID(r.Primary.Attributes["project_id"], r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving variable %s", err)
		}
	}
	return nil
}

func testOctopusDeployVariableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyVarHelper(s, client)
}

func destroyVarHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Variable.DeleteSingle(r.Primary.Attributes["project_id"], r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving variable %s", err)
		}
		return fmt.Errorf("Variable still exists")
	}
	return nil
}
