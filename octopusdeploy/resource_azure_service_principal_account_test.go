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

func TestAccOctopusDeployAzureServicePrincipalAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_principal." + localName

	applicationID := uuid.New()
	applicationPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()
	tenantedDeploymentMode := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	tenantID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "application_password", applicationPassword),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
				),
				Config: testAzureServicePrincipalAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID, applicationPassword, tenantedDeploymentMode),
			},
		},
	})
}

func testAzureServicePrincipalAccountBasic(localName string, name string, description string, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, applicationPassword string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_service_principal" "%s" {
		application_id = "%s"
		application_password = "%s"
		description = "%s"
		name = "%s"
		subscription_id = "%s"
		tenant_id = "%s"
		tenanted_deployment_participation = "%s"
	}`, localName, applicationID, applicationPassword, description, name, subscriptionID, tenantID, tenantedDeploymentParticipation)
}

func testAzureServicePrincipalAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}
