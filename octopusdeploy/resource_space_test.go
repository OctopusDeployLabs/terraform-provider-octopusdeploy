package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSpaceBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeploySpace + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testSpaceDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, constID),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
				Config: testSpaceBasic(localName, name),
			},
		},
	})
}

func testSpaceBasic(localName string, name string) string {
	userLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	userDisplayName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	userPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userUsername := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	return fmt.Sprintf(testUserBasic(userLocalName, userDisplayName, userPassword, userUsername)+"\n"+
		`resource "%s" "%s" {
			name = "%s"
			space_managers_team_members = ["${%s.%s.id}"]
		}`, constOctopusDeploySpace, localName, name, constOctopusDeployUser, userLocalName)
}

func testSpaceExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		spaceID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Spaces.GetByID(spaceID); err != nil {
			return err
		}

		return nil
	}
}

func testSpaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
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
