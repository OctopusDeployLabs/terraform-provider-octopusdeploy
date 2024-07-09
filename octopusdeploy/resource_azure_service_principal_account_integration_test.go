package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"
)

// TestAzureAccountResource verifies that an Azure account can be reimported with the correct settings
func TestAzureAccountResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "4-azureaccount", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := accounts.AccountsQuery{
		PartialName: "Azure",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Accounts.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an account called \"Azure\"")
	}
	resource := resources.Items[0].(*accounts.AzureServicePrincipalAccount)

	if fmt.Sprint(resource.SubscriptionID) != "95bf77d2-64b1-4ed2-9de1-b5451e3881f5" {
		t.Fatalf("The account must be have a client ID of \"95bf77d2-64b1-4ed2-9de1-b5451e3881f5\"")
	}

	if fmt.Sprint(resource.TenantID) != "18eb006b-c3c8-4a72-93cd-fe4b293f82ee" {
		t.Fatalf("The account must be have a client ID of \"18eb006b-c3c8-4a72-93cd-fe4b293f82ee\"")
	}

	if resource.Description != "Azure Account" {
		t.Fatalf("The account must be have a description of \"Azure Account\"")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatalf("The account must be have a tenanted deployment participation of \"Untenanted\"")
	}
}
