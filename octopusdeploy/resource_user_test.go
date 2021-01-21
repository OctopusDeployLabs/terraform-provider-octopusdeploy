package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasic(localName, displayName, true, false, password, username, emailAddress),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUserBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	isActive := true
	isService := false
	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "email_address", emailAddress),
					resource.TestCheckResourceAttr(resourceName, "is_active", strconv.FormatBool(isActive)),
					resource.TestCheckResourceAttr(resourceName, "is_service", strconv.FormatBool(isService)),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
				Config: testAccUserBasic(localName, displayName, isActive, isService, password, username, emailAddress),
			},
			{
				Config:                  testAccUserImport(localName, username),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserImport(localName string, username string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {}`, localName)
}

func testAccUserBasic(localName string, displayName string, isActive bool, isService bool, password string, username string, emailAddress string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {
		display_name  = "%s"
		email_address = "%s"
		is_active     = %v
		is_service    = %v
		password      = "%s"
		username      = "%s"

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
	}`, localName, displayName, emailAddress, isActive, isService, password, username, emailAddress, displayName)
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

func testAccUserCheckDestroy(s *terraform.State) error {
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
