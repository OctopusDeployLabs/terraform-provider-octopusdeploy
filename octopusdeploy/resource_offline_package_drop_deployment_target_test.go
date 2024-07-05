package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

// TestOfflineDropTargetResource verifies that an offline drop can be reimported with the correct settings
func TestOfflineDropTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "34-offlinedroptarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "34a-offlinedroptargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "offline" {
			t.Fatal("The machine must have a role of \"offline\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.OfflinePackageDropEndpoint).ApplicationsDirectory != "c:\\temp" {
			t.Fatal("The machine must have a Endpoint.ApplicationsDirectory of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.OfflinePackageDropEndpoint).ApplicationsDirectory + "\")")
		}

		if resource.Endpoint.(*machines.OfflinePackageDropEndpoint).WorkingDirectory != "c:\\temp" {
			t.Fatal("The machine must have a Endpoint.OctopusWorkingDirectory of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.OfflinePackageDropEndpoint).WorkingDirectory + "\")")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "34a-offlinedroptargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
