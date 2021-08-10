package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployProjectGroupBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testProjectGroupDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccProjectGroupBasic(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project_group." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectGroupDestroy,
		Steps: []resource.TestStep{
			// create projectgroup with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccProjectGroupBasic(localName, name),
			},
			// create update it with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccProjectGroupWithDescription(localName, name, description),
			},
			// update again by remove its description
			{
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccProjectGroupBasic(localName, name),
			},
		},
	})
}

func testAccProjectGroupBasic(localName string, name string) string {
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testAccEnvironment(environmentLocalName, environmentName) + "\n" +
		testAccProjectGroup(localName, name))
}

func testAccProjectGroup(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccProjectGroupWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		description = "%s"
		name        = "%s"
	}`, localName, name, description)
}

func testProjectGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		projectGroupID := rs.Primary.ID
		projectGroup, err := client.ProjectGroups.GetByID(projectGroupID)
		if err == nil {
			if projectGroup != nil {
				return fmt.Errorf("project group (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testProjectGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if _, err := client.ProjectGroups.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := client.ProjectGroups.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}
