package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"
)

// TestGcpAccountResource verifies that a GCP account can be reimported with the correct settings
func TestGcpAccountResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "6-gcpaccount", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := accounts.AccountsQuery{
		PartialName: "Google",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Accounts.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an account called \"Google\"")
	}
	resource := resources.Items[0].(*accounts.GoogleCloudPlatformAccount)

	if !resource.JsonKey.HasValue {
		t.Fatalf("The account must be have a JSON key")
	}

	if resource.Description != "A test account" {
		t.Fatalf("The account must be have a description of \"A test account\"")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatalf("The account must be have a tenanted deployment participation of \"Untenanted\"")
	}

	if len(resource.TenantTags) != 0 {
		t.Fatalf("The account must be have no tenant tags")
	}
}
