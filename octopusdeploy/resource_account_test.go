package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_account." + localName

	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountType := octopusdeploy.AccountTypeAzureServicePrincipal
	applicationID := uuid.New()
	clientSecret := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secretKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
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
					testAccAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "account_type", string(accountType)),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "client_secret", clientSecret),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
				),
				Config: testAccountBasic(localName, name, accountType, applicationID, tenantID, subscriptionID, clientSecret, accessKey, secretKey, tenantedDeploymentMode),
			},
		},
	})
}

func TestAccOctopusCloudTest(t *testing.T) {
	localName := "principal"
	resourceName := "octopusdeploy_account." + localName

	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountType := octopusdeploy.AccountTypeAzureServicePrincipal
	applicationID := uuid.New()
	clientSecret := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := "Principal"
	secretKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
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
					testAccAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "account_type", string(accountType)),
					resource.TestCheckResourceAttr(resourceName, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(resourceName, "client_secret", clientSecret),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
				),
				Config: testAccountBasic(localName, name, accountType, applicationID, tenantID, subscriptionID, clientSecret, accessKey, secretKey, tenantedDeploymentMode),
			},
		},
	})
}

func testAccountBasic(localName string, name string, accountType octopusdeploy.AccountType, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, clientSecret string, accessKey string, secretKey string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_account" "%s" {
		account_type = "%s"
		application_id = "%s"
		client_secret = "%s"
		name = "%s"
		subscription_id = "%s"
		tenant_id = "%s"
		tenanted_deployment_participation = "%s"
	}`, localName, accountType, applicationID, clientSecret, name, subscriptionID, tenantID, tenantedDeploymentParticipation)
}

func testAccAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAccountCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
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
