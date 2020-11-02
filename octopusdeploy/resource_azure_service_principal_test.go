package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployAzureServicePrincipalBasic(t *testing.T) {
	name := acctest.RandString(10)
	clientID := uuid.New().String()
	key := acctest.RandString(10)
	subscriptionID := uuid.New().String()
	tagName := acctest.RandString(10)
	tagSetName := acctest.RandString(10)
	tenantID := uuid.New().String()

	const accountPrefix = constOctopusDeployAzureServicePrincipal + ".foo"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAzureServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureServicePrincipalBasic(tagSetName, tagName, name, clientID, tenantID, subscriptionID, key, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployAzureServicePrincipalExists(accountPrefix),
					resource.TestCheckResourceAttr(accountPrefix, constName, name),
					resource.TestCheckResourceAttr(accountPrefix, constClientID, clientID),
					resource.TestCheckResourceAttr(accountPrefix, constTenantID, tenantID),
					resource.TestCheckResourceAttr(accountPrefix, constSubscriptionNumber, subscriptionID),
					resource.TestCheckResourceAttr(accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(accountPrefix, constTenantedDeploymentParticipation, string(tenantedDeploymentParticipation)),
				),
			},
		},
	})
}

func testAzureServicePrincipalBasic(tagSetName string, tagName string, name string, clientID string, tenantID string, subscriptionID string, clientSecret string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`
		resource "%s" "testtagset" {
			name = "%s"
			tag {
				name = "%s"
				color = "#6e6e6f"
			}
		}

		resource "%s" "foo" {
			name = "%s"
			client_id = "%s"
			tenant_id = "%s"
			subscription_number = "%s"
			key = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		constOctopusDeployTagSet, tagSetName, tagName, constOctopusDeployAzureServicePrincipal, name, clientID, tenantID, subscriptionID, clientSecret, tagSetName, tenantedDeploymentParticipation,
	)
}

func testOctopusDeployAzureServicePrincipalExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsAzureServicePrincipalHelper(s, client)
	}
}

func existsAzureServicePrincipalHelper(s *terraform.State, client *octopusdeploy.Client) error {
	accountID := s.RootModule().Resources[constOctopusDeployAzureServicePrincipal+".foo"].Primary.ID
	if _, err := client.Accounts.GetByID(accountID); err != nil {
		return err
	}

	return nil
}

func testOctopusDeployAzureServicePrincipalDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constOctopusDeployAzureServicePrincipal {
			continue
		}

		accountID := rs.Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
