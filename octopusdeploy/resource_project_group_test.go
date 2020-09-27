package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
						terraformNamePrefix, constName, projectGroupName),
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
						terraformNamePrefix, constName, projectGroupName),
				),
			},
			// create update it with a description
			{
				Config: testAccProjectGroupWithDescription(projectGroupName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, projectGroupName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constDescription, description),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectGroupBasic(projectGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectGroupExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, projectGroupName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constDescription, ""),
				),
			},
		},
	})
}

func testAccProjectGroupBasic(name string) string {
	return fmt.Sprintf(`
		resource constOctopusDeployProjectGroup "foo" {
			name           = "%s"
		  }
		`,
		name,
	)
}
func testAccProjectGroupWithDescription(name, description string) string {
	return fmt.Sprintf(`
		resource constOctopusDeployProjectGroup "foo" {
			name           = "%s"
			description    = "%s"
		  }
		`,
		name, description,
	)
}

func testAccCheckOctopusDeployProjectGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	if err := destroyHelperProjectGroup(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelperProjectGroup(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelperProjectGroup(s *terraform.State, apiClient *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := apiClient.ProjectGroups.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving projectgroup %s", err)
		}
		return fmt.Errorf("projectgroup still exists")
	}
	return nil
}

func existsHelperProjectGroup(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.ProjectGroups.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("received an error retrieving projectgroup %s", err)
		}
	}
	return nil
}
