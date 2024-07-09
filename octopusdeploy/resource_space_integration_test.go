package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"

	"k8s.io/utils/strings/slices"
)

// TestSpaceResource verifies that a space can be reimported with the correct settings
func TestSpaceResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "1-singlespace", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, "", test.ApiKey)
	query := spaces.SpacesQuery{
		IDs:  []string{newSpaceId},
		Skip: 0,
		Take: 1,
	}
	spaces, err := client.Spaces.Get(query)
	space := spaces.Items[0]

	if err != nil {
		t.Fatal(err.Error())
	}

	if space.Description != "My test space" {
		t.Fatalf("New space must have the name \"My test space\"")
	}

	if space.IsDefault {
		t.Fatalf("New space must not be the default one")
	}

	if space.TaskQueueStopped {
		t.Fatalf("New space must not have the task queue stopped")
	}

	if slices.Index(space.SpaceManagersTeams, "teams-administrators") == -1 {
		t.Fatalf("New space must have teams-administrators as a manager team")
	}
}
