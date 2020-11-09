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
	prefix := constOctopusDeployAccount + "." + localName

	accessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountType := octopusdeploy.AccountTypeAzureServicePrincipal
	clientID := uuid.New()
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
					resource.TestCheckResourceAttr(prefix, constAccountType, string(accountType)),
					resource.TestCheckResourceAttr(prefix, constClientID, clientID.String()),
					resource.TestCheckResourceAttr(prefix, constClientSecret, clientSecret),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, constSubscriptionID, subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, constTenantID, tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
				),
				Config: testAccountBasic(localName, name, accountType, clientID, tenantID, subscriptionID, clientSecret, accessKey, secretKey, tenantedDeploymentMode),
			},
		},
	})
}

func testAccountBasic(localName string, name string, accountType octopusdeploy.AccountType, clientID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, clientSecret string, accessKey string, secretKey string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		account_type = "%s"
		client_id = "%s"
		client_secret = "%s"
		name = "%s"
		subscription_id = "%s"
		tenant_id = "%s"
		tenant_tags = ["Hosted Instances/Active"]
		tenanted_deployment_participation = "%s"
	}`, constOctopusDeployAccount, localName, accountType, clientID, clientSecret, name, subscriptionID, tenantID, tenantedDeploymentParticipation)
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
		accountID := rs.Primary.ID
		account, err := client.Accounts.GetByID(accountID)
		if err == nil {
			if account != nil {
				return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
