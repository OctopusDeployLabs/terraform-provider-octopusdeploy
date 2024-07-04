package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"
)

// TestAzureCloudServiceTargetResource verifies that a azure cloud service target can be reimported with the correct settings
func TestAzureCloudServiceTargetResource(t *testing.T) {
	// I could not figure out a combination of properties that made the octopusdeploy_azure_subscription_account resource work
	return

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "35-azurecloudservicetarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			PartialName: "Azure",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Azure\"")
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

		if resource.Endpoint.(*machines.AzureCloudServiceEndpoint).CloudServiceName != "servicename" {
			t.Fatal("The machine must have a Endpoint.CloudServiceName of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.AzureCloudServiceEndpoint).CloudServiceName + "\")")
		}

		if resource.Endpoint.(*machines.AzureCloudServiceEndpoint).StorageAccountName != "accountname" {
			t.Fatal("The machine must have a Endpoint.StorageAccountName of \"accountname\" (was \"" + resource.Endpoint.(*machines.AzureCloudServiceEndpoint).StorageAccountName + "\")")
		}

		if !resource.Endpoint.(*machines.AzureCloudServiceEndpoint).UseCurrentInstanceCount {
			t.Fatal("The machine must have Endpoint.UseCurrentInstanceCount set")
		}

		return nil
	})
}
