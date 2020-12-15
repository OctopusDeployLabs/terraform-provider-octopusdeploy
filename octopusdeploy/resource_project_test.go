package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployProjectBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		CheckDestroy: testAccProjectCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(localName, name, description),
			},
		},
	})
}

func TestAccOctopusDeployProjectWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccProjectCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.iis_website.0.step_name"),
				),
				Config: testAccProjectBasic(localName, name, description),
			},
		},
	})
}

func testAccProjectBasic(localName string, name string, description string) string {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	lifecycleID := "octopusdeploy_lifecycle." + lifecycleLocalName + ".id"
	projectGroupID := "octopusdeploy_project_group." + projectGroupLocalName + ".id"

	return fmt.Sprintf(testAccLifecycleBasic(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroupBasic(projectGroupLocalName, projectGroupName)+"\n"+
		`resource "octopusdeploy_project" "%s" {
		  description      = "%s"
		  lifecycle_id     = %s
		  name             = "%s"
		  project_group_id = %s

		  connectivity_policy {
		    allow_deployments_to_no_targets = true
			skip_machine_behavior           = "None"
		  }
		}
		
		data "octopusdeploy_projects" "%s" {
		  ids = [octopusdeploy_project.%s.id]
		}`, localName, description, lifecycleID, name, projectGroupID, localName, localName)
}

func testAccProjectCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		id := rs.Primary.ID
		switch rs.Type {
		case "octopusdeploy_lifecycle":
			lifecycle, err := client.Lifecycles.GetByID(id)
			if err == nil && lifecycle != nil {
				return fmt.Errorf("lifecycle (%s) still exists", id)
			}
		case "octopusdeploy_project_group":
			projectGroup, err := client.ProjectGroups.GetByID(id)
			if err == nil && projectGroup != nil {
				return fmt.Errorf("project group (%s) still exists", id)
			}
		case "octopusdeploy_project":
			project, err := client.Projects.GetByID(id)
			if err == nil && project != nil {
				return fmt.Errorf("project (%s) still exists", id)
			}
		}
	}

	return nil
}

func testAccCheckOctopusDeployProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyProjectHelper(s, client); err != nil {
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

func destroyProjectHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("error retrieving project %s", err)
		}
		return fmt.Errorf("project still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_project" {
			if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
