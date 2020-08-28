package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployAccountBasic(t *testing.T) {
	const accountPrefix = "octopusdeploy_account.foo"
	const accountName = "Testing one two three"
	const accountType = "Token"
	const tagSetName = "TagSet"
	const tagName = "Tag"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccountBasic(tagSetName, tagName, accountName, accountType),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployAccountExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "name", accountName),
					resource.TestCheckResourceAttr(
						accountPrefix, "account_type", accountType),
				),
			},
		},
	})
}

func testAccountBasic(tagSetName string, tagName string, accountName string, accountType string) string {
	return fmt.Sprintf(`

		resource "octopusdeploy_tag_set" "testtagset" {
			name = "%s"

			tag {
				name = "%s"
				color = "#6e6e6f"
			}
		}


		resource "octopusdeploy_account" "foo" {
			name           = "%s"
			account_type    = "%s"
		}
		`,
		tagSetName, tagName, accountName, accountType,
	)
}

func testOctopusDeployAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsaccountHelper(s, client)
	}
}

func existsaccountHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_account.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving account %s", err)
	}

	return nil
}

func testOctopusDeployAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyaccountHelper(s, client)
}

func destroyaccountHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_account.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		if err == octopusdeploy.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving account %s", err)
	}
	return fmt.Errorf("Account still exists")
}
