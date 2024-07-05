package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

// TestAzureServiceFabricTargetResource verifies that a service fabric target can be reimported with the correct settings
func TestAzureServiceFabricTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "36-servicefabrictarget", []string{
			"-var=target_service_fabric=whatever",
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("..", "terraform", "36a-servicefabrictargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			PartialName: "Service Fabric",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Service Fabric\"")
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

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).ConnectionEndpoint != "http://endpoint" {
			t.Fatal("The machine must have a Endpoint.ConnectionEndpoint of \"http://endpoint\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).ConnectionEndpoint + "\")")
		}

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadCredentialType != "UserCredential" {
			t.Fatal("The machine must have a Endpoint.AadCredentialType of \"UserCredential\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadCredentialType + "\")")
		}

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadUserCredentialUsername != "username" {
			t.Fatal("The machine must have a Endpoint.AadUserCredentialUsername of \"username\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadUserCredentialUsername + "\")")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "36a-servicefabrictargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
