package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserRoleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserRoleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserRoleMinimum(localName, name),
			},
		},
	})
}

func TestAccUserRolePermissions(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserRoleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserRolePermissions(localName, name),
			},
		},
	})
}

func testUserRoleMinimum(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%s" {
		name = "%s"
	}`, localName, name)
}

func testUserRolePermissions(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%s" {
		name                       = "%s"
		granted_space_permissions  = ["AccountCreate"]
		granted_system_permissions = ["SpaceView"]
	}`, localName, name)
}

func testAccUserRoleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_user_role" {
			continue
		}

		_, err := octoClient.UserRoles.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("user role (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
