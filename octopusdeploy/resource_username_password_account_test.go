package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUsernamePasswordBasic(t *testing.T) {
	name := acctest.RandString(10)
	resourceName := constOctopusDeployUsernamePasswordAccount + "." + name

	password := acctest.RandString(10)
	username := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployUsernamePasswordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsernamePasswordBasic(name, username, password),
				Check: resource.ComposeTestCheckFunc(
					testUsernamePasswordExists(name),
					resource.TestCheckResourceAttr(resourceName, constName, name),
					resource.TestCheckResourceAttr(resourceName, constPassword, password),
					resource.TestCheckResourceAttr(resourceName, constUsername, username),
				),
			},
		},
	})
}

func testUsernamePasswordBasic(name string, username string, password string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
			password = "%s"
			username = "%s"
		}
		`,
		constOctopusDeployUsernamePasswordAccount, name, name, password, username,
	)
}

func testUsernamePasswordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		accountID := s.RootModule().Resources[constOctopusDeployUsernamePasswordAccount+"."+name].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}

func testOctopusDeployUsernamePasswordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constOctopusDeployUsernamePasswordAccount {
			continue
		}

		accountID := rs.Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}
		return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
	}

	return nil
}
