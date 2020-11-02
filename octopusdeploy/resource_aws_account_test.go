package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAWSAccountBasic(t *testing.T) {
	accessKey := acctest.RandString(10)
	name := acctest.RandString(10)
	secretKey := acctest.RandString(10)

	const accountPrefix = constOctopusDeployAWSAccount + ".foo"
	const tenantedDeploymentParticipation = octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	var account octopusdeploy.AmazonWebServicesAccount

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAWSAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAWSAccountBasic(name, accessKey, secretKey, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAWSAccountExists(name, &account),
					resource.TestCheckResourceAttr(accountPrefix, constName, name),
					resource.TestCheckResourceAttr(accountPrefix, constAccessKey, accessKey),
					resource.TestCheckResourceAttr(accountPrefix, constSecretKey, secretKey),
					resource.TestCheckResourceAttr(accountPrefix, constTenantedDeploymentParticipation, string(tenantedDeploymentParticipation)),
				),
			},
		},
	})
}

func testAWSAccountBasic(name string, accessKey string, secretKey string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "foo" {
		name = "%s"
		access_key = "%s"
		secret_key = "%s"
		tenanted_deployment_participation = "%s"
	}`, constOctopusDeployAWSAccount, name, accessKey, secretKey, tenantedDeploymentParticipation)
}

func testAWSAccountExists(accountName string, account *octopusdeploy.AmazonWebServicesAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsAWSAccountHelper(s, client)
	}
}

func existsAWSAccountHelper(s *terraform.State, client *octopusdeploy.Client) error {
	// client := testAccProvider.Meta().(*octopusdeploy.Client)
	// query := octopusdeploy.AccountsQuery{PartialName: accountName, Take: 1}
	// accounts, err := client.Accounts.Get(query)
	// if err != nil {
	// 	return err
	// }
	// if len(accounts.Items) != 1 {
	// 	return fmt.Errorf("account not found")
	// }

	// *account = *(accounts.Items[0]).(*octopusdeploy.AmazonWebServicesAccount)
	// return nil

	accountID := s.RootModule().Resources[constOctopusDeployAWSAccount+".foo"].Primary.ID
	if _, err := client.Accounts.GetByID(accountID); err != nil {
		return err
	}

	return nil
}

func testOctopusDeployAWSAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constOctopusDeployAWSAccount {
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
