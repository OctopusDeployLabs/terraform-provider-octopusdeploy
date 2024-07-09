package octopusdeploy

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccScopedUserRole(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	teamResource := "octopusdeploy_team." + localName
	environmentResource := "octopusdeploy_environment." + localName
	userRoleResource := "octopusdeploy_user_role." + localName
	resourceName := "octopusdeploy_scoped_user_role." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	// TODO: replace with client reference
	spaceID := os.Getenv("OCTOPUS_SPACE")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccScopedUserRoleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeAggregateTestCheckFunc(
					testScopedUserRoleExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "user_role_id", userRoleResource, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_id", teamResource, "id"),
					resource.TestCheckResourceAttr(resourceName, "environment_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_ids.0", environmentResource, "id"),
				),
				Config: testAccScopedUserRole(spaceID, localName, name, description),
			},
		},
	})
}

func testScopedUserRoleExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		if _, err := octoClient.ScopedUserRoles.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccScopedUserRoleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_scoped_user_role" {
			continue
		}

		_, err := octoClient.ScopedUserRoles.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("scoped user role (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccScopedUserRole(spaceID string, localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%[2]s" {
		name = "%[4]s"
		granted_space_permissions = ["AccountCreate"]
	}

	resource "octopusdeploy_team" "%[2]s" {
		description = "%[3]s"
		name = "%[4]s"
	}

	resource "octopusdeploy_environment" "%[2]s" {
		name = "%[4]s"
	}
	
	resource "octopusdeploy_scoped_user_role" "%[2]s" {
		space_id = "%[1]s"
		team_id = octopusdeploy_team.%[2]s.id
		user_role_id = octopusdeploy_user_role.%[2]s.id
		environment_ids = [octopusdeploy_environment.%[2]s.id]
	}`, spaceID, localName, description, name)
}
