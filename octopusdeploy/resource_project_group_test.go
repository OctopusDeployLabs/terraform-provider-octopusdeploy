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
	prefix := "octopusdeploy_project_group." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_group." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectGroupDestroy,
		Steps: []resource.TestStep{
			// create projectgroup with no description
			{
				Config: testAccProjectGroupBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
			// create update it with a description
			{
				Config: testAccProjectGroupWithDescription(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectGroupBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", ""),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
		},
	})
}

func testAccProjectGroupBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccProjectGroupWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name        = "%s"
		description = "%s"
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

func testProjectGroupExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		projectGroupID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.ProjectGroups.GetByID(projectGroupID); err != nil {
			return err
		}

		return nil
	}
}

func destroyHelperProjectGroup(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.ProjectGroups.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("error retrieving projectgroup %s", err)
		}
		return fmt.Errorf("projectgroup still exists")
	}
	return nil
}
