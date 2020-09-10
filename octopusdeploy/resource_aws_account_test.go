package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAWSAccountBasic(t *testing.T) {
	const accountPrefix = "octopusdeploy_aws_account.foo"
	const name = "awsaccount"
	const accessKey = "AKIA6DEJDS6OY7FC3I50"
	const secretKey = "x81L4H3riyiWRuBEPlz1"

	const tagSetName = "TagSet"
	const tagName = "Tag"
	var tenantTags = fmt.Sprintf("%s/%s", tagSetName, tagName)
	const tenantedDeploymentParticipation = enum.TenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployAzureServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAWSAccountBasic(tagSetName, tagName, name, accessKey, secretKey, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAWSAccountExists(accountPrefix),
					resource.TestCheckResourceAttr(
						accountPrefix, "name", name),
					resource.TestCheckResourceAttr(
						accountPrefix, "access_key", accessKey),
					resource.TestCheckResourceAttr(
						accountPrefix, "secret_key", secretKey),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenant_tags.0", tenantTags),
					resource.TestCheckResourceAttr(
						accountPrefix, "tenanted_deployment_participation", tenantedDeploymentParticipation.String()),
				),
			},
		},
	})
}

func testIsAccountTypeAWS(t *testing.T) *model.Account {
	accountName := "awsaccounttest"
	testType, err := model.NewAccount(accountName, enum.AmazonWebServicesAccount)

	if err != nil {
		assert.FailNow(t, "The test has failed due to: ", err)
	}

	assert.Error(t, err)
	assert.NotNil(t, err)

	return testType
}

func testIsAWSAccountNil(t *testing.T) *model.Account {
	accountName := "awsaccounttest"
	testNil, err := model.NewAccount(accountName, enum.AmazonWebServicesAccount)

	var pnt *model.Account
	fmt.Printf("Type Account is nil: %v", pnt == nil)

	assert.NotNil(t, err)

	return testNil
}

func testAWSAccountBasic(tagSetName string, tagName string, name string, accessKey string, secretKey string, tenantedDeploymentParticipation enum.TenantedDeploymentMode) string {
	return fmt.Sprintf(`


		resource "octopusdeploy_azure_service_principal" "foo" {
			name           = "%s"
			access_key = "%s"
			secret_key = "%s"
			tagSetName = "%s"
			tenant_tags = ["${octopusdeploy_tag_set.testtagset.name}/%s"]
			tenanted_deployment_participation = "%s"
		}
		`,
		tagSetName, tagName, name, accessKey, secretKey, tenantedDeploymentParticipation,
	)
}

func testAWSAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		return existsAzureServicePrincipalHelper(s, client)
	}
}

func existsAWSAccountHelper(s *terraform.State, client *client.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := client.Accounts.Get(accountID); err != nil {
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}

	return nil
}

func testOctopusDeployAWSAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	return destroyAzureServicePrincipalHelper(s, client)
}

func destroyAWSAccountHelper(s *terraform.State, apiClient *client.Client) error {

	accountID := s.RootModule().Resources["octopusdeploy_azure_service_principal.foo"].Primary.ID

	if _, err := apiClient.Accounts.Get(accountID); err != nil {
		if err == client.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving azure service principal %s", err)
	}
	return fmt.Errorf("Azure Service Principal still exists")
}
