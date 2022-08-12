package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccountCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_account" {
			continue
		}

		account, err := client.Accounts.GetByID(rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
