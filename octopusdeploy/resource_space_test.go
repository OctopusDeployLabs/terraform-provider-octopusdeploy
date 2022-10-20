package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSpaceImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccSpaceCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSpaceBasic(localName, name, slug),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpaceBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_space." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccSpaceCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "space_managers_teams.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "space_managers_teams.0"),
				),
				Config: testSpaceBasic(localName, name, slug),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "space_managers_teams.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "space_managers_teams.0"),
				),
				Config: testSpaceDataSource(localName, newName, slug),
			},
		},
	})
}

func testSpaceDataSource(localName string, name string, slug string) string {
	return fmt.Sprintf(testSpaceBasic(localName, name, slug)+"\n"+
		`data "octopusdeploy_spaces" "%s" {
			partial_name = "%s"
		}`, localName, name)
}

func testSpaceBasic(localName string, name string, slug string) string {
	userLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userDisplayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userEmailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	userPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testAccUserBasic(userLocalName, userDisplayName, true, false, userPassword, userUsername, userEmailAddress)+"\n"+
		`resource "octopusdeploy_space" "%s" {
			name = "%s"
			slug = "%s"
			space_managers_teams  = ["teams-managers"]

			lifecycle {
			  ignore_changes = [space_managers_teams]
			}
		}`, localName, name, slug)
}

func testSpaceExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		spaceID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Spaces.GetByID(spaceID); err != nil {
			return err
		}

		return nil
	}
}

func testAccSpaceCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		spaceID := rs.Primary.ID
		space, err := client.Spaces.GetByID(spaceID)
		if err == nil {
			if space != nil {
				return fmt.Errorf("space (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
