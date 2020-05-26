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
	const accountType = "AzureServicePrincipal"
	const clientID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e1"
	const tenantID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e2"
	const subscriptionID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e3"
	//nolint:gosec
	const clientSecret = "18eb006b-c3c8-4a72-93cd-fe4b293f82e4"
	const tagSetName = "TagSet"
	const tagName = "Tag"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = "TenantedOrUntenanted"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccountBasic(tagSetName, tagName, accountName, accountType, clientID, tenantID, subscriptionID, clientSecret, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployAccountExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "name", accountName),
					resource.TestCheckResourceAttr(
						accountPrefix, "account_type", accountType),
					resource.TestCheckResourceAttr(
						accountPrefix, "client_id", clientID),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_id", tenantID),
					resource.TestCheckResourceAttr(
						accountPrefix, "subscription_id", subscriptionID),
					resource.TestCheckResourceAttr(
						accountPrefix, "client_secret", clientSecret),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenanted_deployment_participation", tenantedDeploymentParticipation),
				),
			},
		},
	})
}

func testAccountBasic(tagSetName string, tagName string, accountName string, accountType, clientID string, tenantID string, subscriptionID string, clientSecret string, tenantedDeploymentParticipation string) string {
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
			client_id = "%s"
			tenant_id = "%s"
			subscription_id = "%s"
			client_secret = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		tagSetName, tagName, accountName, accountType, clientID, tenantID, subscriptionID, clientSecret, tagName, tenantedDeploymentParticipation,
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
