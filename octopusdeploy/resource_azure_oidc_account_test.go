package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOctopusDeployAzureOpenIDConnectAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_openid_connect." + localName

	applicationID := uuid.New()
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()
	tenantedDeploymentMode := core.TenantedDeploymentModeTenantedOrUntenanted
	tenantID := uuid.New()

	executionKeys := []string{"space"}
	healthKeys := []string{"target"}
	accountKeys := []string{"type"}
	audience := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccountCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys", healthKeys[0]),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys", accountKeys[0]),
					resource.TestCheckResourceAttr(prefix, "audience", audience),
				),
				Config: testAzureOpenIDConnectAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID, tenantedDeploymentMode, executionKeys, healthKeys, accountKeys, audience),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys", healthKeys[0]),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys", accountKeys[0]),
					resource.TestCheckResourceAttr(prefix, "audience", audience),
				),
				Config: testAzureOpenIDConnectAccountBasic(localName, name, newDescription, applicationID, tenantID, subscriptionID, tenantedDeploymentMode, executionKeys, healthKeys, accountKeys, audience),
			},
		},
	})
}

func testAzureOpenIDConnectAccountBasic(localName string, name string, description string, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, tenantedDeploymentParticipation core.TenantedDeploymentMode, execution_subject_keys []string, health_subject_keys []string, account_test_subject_keys []string, audience string) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_openid_connect" "%s" {
		application_id = "%s"
		description = "%s"
		name = "%s"
		subscription_id = "%s"
		tenant_id = "%s"
		tenanted_deployment_participation = "%s"
		execution_subject_keys = "%s"
		health_subject_keys = "%s"
		account_test_subject_keys = "%s"
		audience = "%s"
	}
	
	data "octopusdeploy_accounts" "test" {
		ids = [octopusdeploy_azure_openid_connect.%s.id]
	}`, localName, applicationID, description, name, subscriptionID, tenantID, tenantedDeploymentParticipation, StringArrayToTerraformArrayFormat(execution_subject_keys), StringArrayToTerraformArrayFormat(health_subject_keys), StringArrayToTerraformArrayFormat(account_test_subject_keys), audience, localName)
}
