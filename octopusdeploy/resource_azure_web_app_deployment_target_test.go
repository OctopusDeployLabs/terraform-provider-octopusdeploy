package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
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
		CheckDestroy: testAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

// TestAzureWebAppTargetResource verifies that a web app target can be reimported with the correct settings
func TestAzureWebAppTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "37-webapptarget", []string{
			"-var=account_sales_account=whatever",
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "37a-webapptarget"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			PartialName: "Web App",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Web App\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).ResourceGroupName != "mattc-webapp" {
			t.Fatal("The machine must have a Endpoint.ResourceGroupName of \"mattc-webapp\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).ResourceGroupName + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppName != "mattc-webapp" {
			t.Fatal("The machine must have a Endpoint.WebAppName of \"mattc-webapp\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppName + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppSlotName != "slot1" {
			t.Fatal("The machine must have a Endpoint.WebAppSlotName of \"slot1\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppSlotName + "\")")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "37a-webapptarget"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
