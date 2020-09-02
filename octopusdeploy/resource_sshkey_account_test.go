package model

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestSSHKeyBasic(t *testing.T) {
	const accountPrefix = "octopusdeploy_sshkey_account.foo"
	const username = "foo"
	const passphrase = "H3lloWorld"

	const tagSetName = "TagSet"
	const tagName = "Tag"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = enum.TenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAzureServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSSHKeyBasic(tagSetName, tagName, username, passphrase, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testSSHKeyExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "username", username),
					resource.TestCheckResourceAttr(
						accountPrefix, "passphrase", passphrase),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenanted_deployment_participation", tenantedDeploymentParticipation.String()),
				),
			},
		},
	})
}

func testSSHKeyBasic(tagSetName string, tagName string, username string, passphrase string, tenantedDeploymentParticipation enum.TenantedDeploymentMode) string {
	return fmt.Sprintf(`


		resource "octopusdeploy_azure_service_principal" "foo" {
			usernamename           = "%s"
			passphrase = "%s"
			tagSetName = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		tagSetName, tagName, username, passphrase, tenantedDeploymentParticipation,
	)
}

func testSSHKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		return existsAzureServicePrincipalHelper(s, client)
	}
}

func existsSSHKeyHelper(s *terraform.State, client *client.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := client.Accounts.Get(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}

	return nil
}

func testOctopusDeploySSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	return destroyAzureServicePrincipalHelper(s, client)
}

func destroySSHKeyHelper(s *terraform.State, apiClient *client.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := apiClient.Accounts.Get(accountID); err != nil {
		if err == client.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}
	return fmt.Errorf("Azure Service Principal still exists")
}
