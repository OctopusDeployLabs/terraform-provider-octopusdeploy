package octopusdeploy

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccUserRoleBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testUserRoleMinimum(localName, name),
			},
		},
	})
}

func (suite *IntegrationTestSuite) TestAccUserRolePermissions() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
