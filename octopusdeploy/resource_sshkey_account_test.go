package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSSHKeyBasic(t *testing.T) {
	const accountPrefix = constOctopusDeploySSHKeyAccount + ".foo"
	const username = "foo"
	const passphrase = "H3lloWorld"

	const tagSetName = "TagSet"
	const tagName = "Tag"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeploySSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSSHKeyBasic(tagSetName, tagName, username, passphrase, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testSSHKeyExists(accountPrefix),
					resource.TestCheckResourceAttr(accountPrefix, constUsername, username),
					resource.TestCheckResourceAttr(accountPrefix, constPassphrase, passphrase),
					resource.TestCheckResourceAttr(accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(accountPrefix, constTenantedDeploymentParticipation, string(tenantedDeploymentParticipation)),
				),
			},
		},
	})
}

func testSSHKeyBasic(tagSetName string, tagName string, username string, passphrase string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`
		resource "%s" "foo" {
			username = "%s"
			passphrase = "%s"
			tagSetName = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		constOctopusDeploySSHKeyAccount, tagSetName, tagName, username, passphrase, tenantedDeploymentParticipation,
	)
}

func testSSHKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsSSHKeyHelper(s, client)
	}
}

func existsSSHKeyHelper(s *terraform.State, client *octopusdeploy.Client) error {
	accountID := s.RootModule().Resources[constOctopusDeploySSHKeyAccount+".foo"].Primary.ID
	if _, err := client.Accounts.GetByID(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving SSH key account %s", err)
	}

	return nil
}

func testOctopusDeploySSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constOctopusDeploySSHKeyAccount {
			continue
		}

		accountID := rs.Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
