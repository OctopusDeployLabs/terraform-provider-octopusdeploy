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
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployAzureServicePrincipal + "." + localName

	name := acctest.RandString(10)
	clientID := uuid.New()
	clientSecret := acctest.RandString(10)
	subscriptionID := uuid.New()
	tenantedDeploymentMode := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	tenantID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAzureServicePrincipalExists(prefix),
					resource.TestCheckResourceAttr(prefix, constClientID, clientID.String()),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constSubscriptionNumber, subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, constTenantID, tenantID.String()),
					resource.TestCheckResourceAttr(prefix, constTenantedDeploymentParticipation, string(tenantedDeploymentMode)),
				),
				Config: testAzureServicePrincipalBasic(localName, name, clientID, tenantID, subscriptionID, clientSecret, tenantedDeploymentMode),
			},
		},
	})
}

func testAzureServicePrincipalBasic(localName string, name string, clientID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, clientSecret string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		client_id = "%s"
		name = "%s"
		subscription_number = "%s"
		tenant_id = "%s"
		key = "%s"
		tenanted_deployment_participation = "%s"
	}`, constOctopusDeployAzureServicePrincipal, localName, clientID, name, subscriptionID, tenantID, clientSecret, tenantedDeploymentParticipation)
}

func testAzureServicePrincipalExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}
