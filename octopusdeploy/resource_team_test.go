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
					resource.TestCheckResourceAttr(resourceName, "description", newDescription),
				),
				Config: testAccTeamBasic(localName, name, newDescription),
			},
		},
	})
}

func testAccTeamBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "%s" {
		description = "%s"
		name = "%s"
	}`, localName, description, name)
}

func testAccTeamCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
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

		if team, err := client.Teams.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("team (%s) still exists", team.GetID())
		}
	}

	return nil
}
