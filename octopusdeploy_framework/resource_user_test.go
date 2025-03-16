package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccUserImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
		CheckDestroy:             testAccUserCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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

func testAccUserImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
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
		userID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := users.GetByID(octoClient, userID); err != nil {
			return err
		}

		return nil
	}
}

func testAccUserCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_user" {
			continue
		}

		_, err := users.GetByID(octoClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("user (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
