package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
)

// TestProjectGroupResource verifies that a project group can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestProjectGroupResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "2-projectgroup", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "2a-projectgroupds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projectgroups.ProjectGroupsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.ProjectGroups.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project group called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test Description" {
		t.Fatalf("The project group must be have a description of \"Test Description\"")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "2a-projectgroupds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
