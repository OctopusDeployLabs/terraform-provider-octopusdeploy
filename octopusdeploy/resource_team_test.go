package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTeamBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_team." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccTeamCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
				Config: testAccTeamBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", newDescription),
				),
				Config: testAccTeamBasic(localName, name, newDescription),
			},
		},
	})
}

func TestAccTeamUserRole(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_team." + localName
	userRoleResource := "octopusdeploy_user_role." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccTeamCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeAggregateTestCheckFunc(
					testTeamExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrPair(resourceName, "user_role.0.user_role_id", userRoleResource, "id"),
					resource.TestCheckResourceAttr(resourceName, "user_role.#", "1"),
				),
				Config: testAccTeamUserRole(localName, name, description),
			},
		},
	})
}

func testTeamExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if _, err := client.Teams.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTeamCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_team" {
			continue
		}

		_, err := client.Teams.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("team (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccTeamBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "%s" {
		description = "%s"
		name = "%s"
	}`, localName, description, name)
}

func testAccTeamUserRole(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%[1]s" {
		name = "%[3]s"
		granted_space_permissions = ["AccountCreate"]
	}

	resource "octopusdeploy_team" "%[1]s" {
		description = "%[2]s"
		name = "%[3]s"

		user_role {
			space_id = "Spaces-1"
			user_role_id = octopusdeploy_user_role.%[1]s.id
		}
	}`, localName, description, name)
}
