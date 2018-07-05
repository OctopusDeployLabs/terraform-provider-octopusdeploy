package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployProjectBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const projectGroupID = "ProjectGroups-1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, projectGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "project_group_id", projectGroupID),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectWithUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const projectGroupID = "ProjectGroups-1"
	const description = "I am a new description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			// create project with no description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, projectGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "project_group_id", projectGroupID),
				),
			},
			// create update it with a description
			{
				Config: testAccProjectWithDescription(projectName, lifeCycleID, projectGroupID, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "project_group_id", projectGroupID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", description),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, projectGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "project_group_id", projectGroupID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", ""),
				),
			},
		},
	})
}

func testAccProjectBasic(name, lifeCycleID, projectGroupID string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project" "foo" {
			name           = "%s"
			lifecycle_id    = "%s"
			project_group_id = "%s"
		  }
		`,
		name, lifeCycleID, projectGroupID,
	)
}
func testAccProjectWithDescription(name, lifeCycleID, projectGroupID, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project" "foo" {
			name           = "%s"
			lifecycle_id    = "%s"
			project_group_id = "%s"
			description    = "%s"
		  }
		`,
		name, lifeCycleID, projectGroupID, description,
	)
}

func testAccCheckOctopusDeployProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Project.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving project %s", err)
		}
		return fmt.Errorf("Project still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Project.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving project %s", err)
		}
	}
	return nil
}
