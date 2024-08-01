package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
)

// TestPollingTargetResource verifies that a polling machine can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestPollingTargetResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "32-pollingtarget", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "32a-pollingtargetds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := machines.MachinesQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Machines.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a machine called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.Endpoint.(*machines.PollingTentacleEndpoint).URI.Host != "abcdefghijklmnopqrst" {
		t.Fatal("The machine must have a Uri of \"poll://abcdefghijklmnopqrst/\" (was \"" + resource.Endpoint.(*machines.PollingTentacleEndpoint).URI.Host + "\")")
	}

	if resource.Thumbprint != "1854A302E5D9EAC1CAA3DA1F5249F82C28BB2B86" {
		t.Fatal("The machine must have a Thumbprint of \"1854A302E5D9EAC1CAA3DA1F5249F82C28BB2B86\" (was \"" + resource.Thumbprint + "\")")
	}

	if len(resource.Roles) != 1 {
		t.Fatal("The machine must have 1 role")
	}

	if resource.Roles[0] != "vm" {
		t.Fatal("The machine must have a role of \"vm\" (was \"" + resource.Roles[0] + "\")")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "32a-pollingtargetds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
