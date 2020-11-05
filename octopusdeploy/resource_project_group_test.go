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
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployProjectGroup + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectGroupWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployProjectGroup + "." + localName

	description := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

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
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
			// create update it with a description
			{
				Config: testAccProjectGroupWithDescription(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, constDescription, description),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectGroupBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(prefix),
					resource.TestCheckResourceAttr(prefix, constDescription, ""),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
		},
	})
}

func testAccProjectGroupBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
	}`, constOctopusDeployProjectGroup, localName, name)
}

func testAccProjectGroupWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name        = "%s"
		description = "%s"
	}`, constOctopusDeployProjectGroup, localName, name, description)
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
			return fmt.Errorf("Received an error retrieving projectgroup %s", err)
		}
		return fmt.Errorf("projectgroup still exists")
	}
	return nil
}
