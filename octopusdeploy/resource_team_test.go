package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTeamBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_team." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
				Config: testAccTeamBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
				Config: testAccTeamBasic(localName, name, updatedDescription),
			},
		},
	})
}

func TestAccTeamUserRole(t *testing.T) {
	SkipCI(t, "error creating user role for team Teams-3: octopus deploy api returned an error on endpoint /api/scopeduserroles - [You cannot use a role with Space level permissions at the System level. Space level permissions: AccountCreate]")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_team." + localName
	userRoleResource := "octopusdeploy_user_role." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	spaceID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTeamCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrPair(resourceName, "user_role.0.user_role_id", userRoleResource, "id"),
					resource.TestCheckResourceAttr(resourceName, "user_role.#", "1"),
				),
				Config: testAccTeamUserRole(spaceID, localName, name, description),
			},
		},
	})
}

func testAccTeamCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := octoClient.Teams.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTeamCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_team" {
			continue
		}

		if team, err := octoClient.Teams.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("team (%s) still exists", team.GetID())
		}
	}

	return nil
}

func testAccTeamBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "%s" {
		description = "%s"
		name        = "%s"
	}`, localName, description, name)
}

func testAccTeamUserRole(spaceID string, localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%[2]s" {
		granted_space_permissions = ["AccountCreate"]
		name                      = "%[4]s"
	}

	resource "octopusdeploy_team" "%[2]s" {
		description = "%[3]s"
		name        = "%[4]s"

		user_role {
			space_id = "%[1]s"
			user_role_id = octopusdeploy_user_role.%[2]s.id
		}
	}`, spaceID, localName, description, name)
}
