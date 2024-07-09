package octopusdeploy

import (
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

// TestGitCredentialsResource verifies that a git credential can be reimported with the correct settings
func TestGitCredentialsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "22-gitcredentialtest", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "22a-gitcredentialtestds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "22a-gitcredentialtestds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup == "" {
		t.Fatal("The target lookup did not succeed.")
	}
}
