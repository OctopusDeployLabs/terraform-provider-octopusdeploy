package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

// TestGitCredentialsResource verifies that a git credential can be reimported with the correct settings
func TestGitCredentialsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "22-gitcredentialtest", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "22a-gitcredentialtestds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "22a-gitcredentialtestds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup == "" {
			t.Fatal("The target lookup did not succeed.")
		}

		return nil
	})
}
