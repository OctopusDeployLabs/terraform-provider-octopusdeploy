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
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	accountPrefix := constOctopusDeployAccount + "." + name

	accessKey := acctest.RandString(10)
	accountType := octopusdeploy.AccountTypeAzureServicePrincipal
	clientID := uuid.New()
	clientSecret := acctest.RandString(10)
	secretKey := acctest.RandString(10)
	subscriptionID := uuid.New()
	tenantedDeploymentMode := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	tenantID := uuid.New()

	var account octopusdeploy.AccountResource

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccountBasic(name, accountType, clientID, tenantID, subscriptionID, clientSecret, accessKey, secretKey, tenantedDeploymentMode),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(name, &account),
					resource.TestCheckResourceAttr(accountPrefix, constName, name),
					resource.TestCheckResourceAttr(accountPrefix, constAccountType, string(accountType)),
					resource.TestCheckResourceAttr(accountPrefix, constClientID, clientID.String()),
					resource.TestCheckResourceAttr(accountPrefix, constTenantID, tenantID.String()),
					resource.TestCheckResourceAttr(accountPrefix, constSubscriptionID, subscriptionID.String()),
					resource.TestCheckResourceAttr(accountPrefix, constClientSecret, clientSecret),
					resource.TestCheckResourceAttr(accountPrefix, constTenantedDeploymentParticipation, string(tenantedDeploymentMode)),
				),
			},
		},
	})
}

func testAccountBasic(name string, accountType octopusdeploy.AccountType, clientID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, clientSecret string, accessKey string, secretKey string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
		account_type = "%s"
		client_id = "%s"
		tenant_id = "%s"
		subscription_id = "%s"
		client_secret = "%s"
		tenant_tags = ["Hosted Instances/Active"]
		tenanted_deployment_participation = "%s"
	}`, constOctopusDeployAccount, name, name, accountType, clientID, tenantID, subscriptionID, clientSecret, tenantedDeploymentParticipation)
}

func testAccountExists(name string, account *octopusdeploy.AccountResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		accountID := s.RootModule().Resources[constOctopusDeployAccount+"."+name].Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}

func testOctopusDeployAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constOctopusDeployAccount {
			continue
		}

		accountID := rs.Primary.ID
		if _, err := client.Accounts.GetByID(accountID); err != nil {
			return err
		}
		return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
	}

	return nil
}
