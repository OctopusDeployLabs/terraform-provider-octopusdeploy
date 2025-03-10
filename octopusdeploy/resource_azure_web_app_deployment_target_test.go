package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOctopusDeployAzureWebAppDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_web_app_deployment_target." + localName
	tenantedDeploymentMode := core.TenantedDeploymentModeTenantedOrUntenanted
	webAppName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccountCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentTargetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
					resource.TestCheckResourceAttr(resourceName, "web_app_name", webAppName),
				),
				Config: testAzureWebAppDeploymentTargetBasic(localName, name, tenantedDeploymentMode, webAppName),
			},
		},
	})
}

func testAzureWebAppDeploymentTargetBasic(localName string, name string, tenantedDeploymentParticipation core.TenantedDeploymentMode, webAppName string) string {
	allowDynamicInfrastructure := false
	azureAccAccountID := uuid.New()
	azureAccLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	azureAccName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	azureAccDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	azureAccPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	azureAccSubscriptionID := uuid.New()
	azureAccTenantedDeploymentMode := core.TenantedDeploymentModeTenantedOrUntenanted
	azureAccTenantID := uuid.New()
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(5, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAzureServicePrincipalAccountBasic(azureAccLocalName, azureAccName, azureAccDescription, azureAccAccountID, azureAccTenantID, azureAccSubscriptionID, azureAccPassword, azureAccTenantedDeploymentMode)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		`resource "octopusdeploy_azure_web_app_deployment_target" "%s" {
			account_id                        = octopusdeploy_azure_service_principal.%s.id
			environments                      = [octopusdeploy_environment.%s.id]
			name                              = "%s"
			resource_group_name               = "%s"
			roles                             = ["test role"]
			tenanted_deployment_participation = "%s"
			web_app_name                      = "%s"
		}`, localName, azureAccLocalName, environmentLocalName, name, resourceGroupName, tenantedDeploymentParticipation, webAppName)
}
