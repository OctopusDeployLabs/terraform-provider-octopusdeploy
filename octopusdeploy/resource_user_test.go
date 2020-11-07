package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestUserBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployUser + "." + localName

	displayName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testUserDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserBasic(localName, displayName, password, username),
				Check: resource.ComposeTestCheckFunc(
					testUserExists(prefix),
					resource.TestCheckResourceAttr(prefix, constDisplayName, displayName),
					resource.TestCheckResourceAttr(prefix, constPassword, password),
					resource.TestCheckResourceAttr(prefix, constUsername, username),
				),
			},
		},
	})
}

func testUserBasic(localName string, displayName string, password string, username string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		display_name = "%s"
		password     = "%s"
		username     = "%s"
		identities   = [
			{
				identity_provider_name = "Octopus ID"
				claims                 = [
					{
						email = {
							value                = "bob.smith@example.com"
							is_identifying_claim = true
						}
						dn    = {
							value                = "Bob Smith"
							is_identifying_claim = false
						}
					}
				]
			}
		]
	}`, constOctopusDeployUser, localName, displayName, password, username)
}

func testUserExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		userID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Users.GetByID(userID); err != nil {
			return err
		}

		return nil
	}
}

func testUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		userID := rs.Primary.ID
		user, err := client.Users.GetByID(userID)
		if err == nil {
			if user != nil {
				return fmt.Errorf("user (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
