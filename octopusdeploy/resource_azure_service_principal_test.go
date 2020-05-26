package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployAzureServicePrinciaplBasic(t *testing.T) {
	const accountPrefix = "octopusdeploy_azure_service_principal.foo"
	const accountName = "Testing one two three"
	const clientID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e1"
	const tenantID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e2"
	const subscriptionID = "18eb006b-c3c8-4a72-93cd-fe4b293f82e3"
	const key = "18eb006b-c3c8-4a72-93cd-fe4b293f82e4"
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
				Config: testAzureServicePrincipalBasic(tagSetName, tagName, accountName, clientID, tenantID, subscriptionID, key, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployAzureServicePrincipalExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "name", accountName),
					resource.TestCheckResourceAttr(
						accountPrefix, "client_id", clientID),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_id", tenantID),
					resource.TestCheckResourceAttr(
						accountPrefix, "subscription_number", subscriptionID),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenanted_deployment_participation", tenantedDeploymentParticipation.String()),
				),
			},
		},
	})
}

func testAzureServicePrincipalBasic(tagSetName string, tagName string, accountName string, clientID string, tenantID string, subscriptionID string, clientSecret string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`

		resource "octopusdeploy_tag_set" "testtagset" {
			name = "%s"

			tag {
				name = "%s"
				color = "#6e6e6f"
			}
		}


		resource "octopusdeploy_azure_service_principal" "foo" {
			name           = "%s"
			client_id = "%s"
			tenant_id = "%s"
			subscription_number = "%s"
			key = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		tagSetName, tagName, accountName, clientID, tenantID, subscriptionID, clientSecret, tagName, tenantedDeploymentParticipation,
	)
}

func testOctopusDeployAzureServicePrincipalExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsAzureServicePrincipalHelper(s, client)
	}
}

func existsAzureServicePrincipalHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}

	return nil
}

func testOctopusDeployAzureServicePrincipalDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyAzureServicePrincipalHelper(s, client)
}

func destroyAzureServicePrincipalHelper(s *terraform.State, client *octopusdeploy.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := client.Account.Get(accountID); err != nil {
		if err == octopusdeploy.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}
	return fmt.Errorf("Azure Service Principal still exists")
}
