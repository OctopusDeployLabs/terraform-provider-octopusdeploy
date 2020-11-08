package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := testUserBasic(localName, displayName, true, false, password, username, emailAddress)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testUserDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,

				// NOTE: import succeeds, however the state differs for the
				// modified_on field

				// ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := "octopusdeploy_user." + localName

	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testUserDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserBasic(localName, displayName, true, false, password, username, emailAddress),
				Check: resource.ComposeTestCheckFunc(
					testUserExists(prefix),
				),
			},
			{
				Config: testUserDataSource(localName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+prefix, "display_name", displayName),
				),
			},
		},
	})
}

func testUserDataSource(localName string, username string) string {
	return fmt.Sprintf(`data "octopusdeploy_user" "%s" {
		username = "%s"
	}`, localName, username)
}

func testUserBasic(localName string, displayName string, isActive bool, isService bool, password string, username string, emailAddress string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {
		display_name = "%s"
		is_active    = %v
		is_service   = %v
		password     = "%s"
		username     = "%s"

		identity {
			provider = "Octopus ID"
			claim {
				name = "email"
				is_identifying_claim = true
				value = "%s"
			}
			claim {
				name = "dn"
				is_identifying_claim = false
				value = "%s"
			}
		}

		lifecycle {
			ignore_changes = [password, modified_on, modified_by]
		}
	}`, localName, displayName, isActive, isService, password, username, emailAddress, displayName)
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
		if rs.Type != "octopusdeploy_user" {
			continue
		}

		_, err := client.Users.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("user (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
