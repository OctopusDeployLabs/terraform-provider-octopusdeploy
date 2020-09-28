package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestTokenAccountBasic(t *testing.T) {
	const accountPrefix = "octopusdeploy_token_account.foo"
	const name = "tokenaccount"
	const token = "SomeRandomTokenValue"

	const tagSetName = "TagSet"
	const tagName = "Tag"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = octopusdeploy.TenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAzureServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTokenAccountBasic(tagSetName, tagName, name, token, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAWSAccountExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "name", name),
					resource.TestCheckResourceAttr(
						accountPrefix, "token", token),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenanted_deployment_participation", tenantedDeploymentParticipation.String()),
				),
			},
		},
	})
}

func testTokenAccountBasic(tagSetName string, tagName string, name string, token string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`


		resource "octopusdeploy_azure_service_principal" "foo" {
			name           = "%s"
			token = "%s"
			tagSetName = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		tagSetName, tagName, name, token, tenantedDeploymentParticipation,
	)
}

func testTokenAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsAzureServicePrincipalHelper(s, client)
	}
}

func existsTokenAccountHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_token_service_principal.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving token service principal %s", err)
	}

	return nil
}

func testOctopusDeployTokenAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyTokenServicePrincipalHelper(s, client)
}

func destroyTokenAccountHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_token_service_principal.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		if err == octopusdeploy.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving token service principal %s", err)
	}
	return fmt.Errorf("Token Service Principal still exists")
}

func destroyTokenServicePrincipalHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_token_service_principal.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		if err == octopusdeploy.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving token service principal %s", err)
	}
	return fmt.Errorf("Token Service Principal still exists")
}
