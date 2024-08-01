package octopusdeploy

import (
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"strings"
)

// TestRunbookResource verifies that a runbook can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestRunbookResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "46-runbooks", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	//err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "46a-runbooks"), newSpaceId, []string{})
	//
	//if err != nil {
	//	return err
	//}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	resources, err := client.Runbooks.GetAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	found := false
	runbookId := ""
	for _, r := range resources {
		if r.Name == "Runbook" {
			found = true
			runbookId = r.ID

			if r.Description != "Test Runbook" {
				t.Fatal("The runbook must be have a description of \"Test Runbook\" (was \"" + r.Description + "\")")
			}

			if r.ConnectivityPolicy.AllowDeploymentsToNoTargets {
				t.Fatal("The runbook must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
			}

			if r.ConnectivityPolicy.ExcludeUnhealthyTargets {
				t.Fatal("The runbook must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
			}

			if r.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
				t.Log("BUG: The runbook must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + r.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
			}

			if r.RunRetentionPolicy.QuantityToKeep != 10 {
				t.Fatal("The runbook must not have RunRetentionPolicy.QuantityToKeep of 10 (was \"" + fmt.Sprint(r.RunRetentionPolicy.QuantityToKeep) + "\")")
			}

			if r.RunRetentionPolicy.ShouldKeepForever {
				t.Fatal("The runbook must not have RunRetentionPolicy.ShouldKeepForever of false (was \"" + fmt.Sprint(r.RunRetentionPolicy.ShouldKeepForever) + "\")")
			}

			if r.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
				t.Log("BUG: The runbook must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + r.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
			}

			if r.MultiTenancyMode != "Untenanted" {
				t.Fatal("The runbook must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + r.MultiTenancyMode + "\")")
			}

			if r.EnvironmentScope != "Specified" {
				t.Fatal("The runbook must be have a EnvironmentScope of \"Specified\" (was \"" + r.EnvironmentScope + "\")")
			}

			if len(r.Environments) != 1 {
				t.Fatal("The runbook must be have a Environments array of 1 (was \"" + strings.Join(r.Environments, ", ") + "\")")
			}

			if r.DefaultGuidedFailureMode != "EnvironmentDefault" {
				t.Fatal("The runbook must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + r.DefaultGuidedFailureMode + "\")")
			}

			if !r.ForcePackageDownload {
				t.Log("BUG: The runbook must be have a ForcePackageDownload of \"true\" (was \"" + fmt.Sprint(r.ForcePackageDownload) + "\")")
			}

			process, err := client.RunbookProcesses.GetByID(r.RunbookProcessID)

			if err != nil {
				t.Fatal("Failed to retrieve the runbook process.")
			}

			if len(process.Steps) != 1 {
				t.Fatal("The runbook must be have a 1 step")
			}
		}
	}

	if !found {
		t.Fatalf("Space must have a runbook called \"Runbook\"")
	}

	// There was an issue where deleting a runbook and reapplying the terraform module caused an error, so
	// verify this process works.
	client.Runbooks.DeleteByID(runbookId)
	err = testFramework.TerraformApply(t, "../terraform/46-runbooks", octoContainer.URI, newSpaceId, []string{})

	if err != nil {
		t.Fatal("Failed to reapply the runbooks after deleting them.")
	}

	// Verify the environment data lookups work
	//lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "46a-runbooks"), "data_lookup")
	//
	//if err != nil {
	//	return err
	//}
	//
	//if lookup != resource.ID {
	//	t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	//}
}
