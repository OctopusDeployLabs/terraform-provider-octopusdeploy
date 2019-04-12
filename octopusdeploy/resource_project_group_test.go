package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployProjectGroupBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project_group.foo"
	const projectGroupName = "Funky Group"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupBasic(projectGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectGroupName),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupWithUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project_group.foo"
	const projectGroupName = "Funky Group"
	const description = "I am a new group description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectGroupDestroy,
		Steps: []resource.TestStep{
			// create projectgroup with no description
			{
				Config: testAccProjectGroupBasic(projectGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectGroupName),
				),
			},
			// create update it with a description
			{
				Config: testAccProjectGroupWithDescription(projectGroupName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectGroupName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", description),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectGroupBasic(projectGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectGroupName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", ""),
				),
			},
		},
	})
}

func testAccProjectGroupBasic(name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name           = "%s"
		  }
		`,
		name,
	)
}
func testAccProjectGroupWithDescription(name, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name           = "%s"
			description    = "%s"
		  }
		`,
		name, description,
	)
}

func testAccCheckOctopusDeployProjectGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyHelperProjectGroup(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelperProjectGroup(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelperProjectGroup(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.ProjectGroup.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving projectgroup %s", err)
		}
		return fmt.Errorf("projectgroup still exists")
	}
	return nil
}

func existsHelperProjectGroup(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.ProjectGroup.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving projectgroup %s", err)
		}
	}
	return nil
}
