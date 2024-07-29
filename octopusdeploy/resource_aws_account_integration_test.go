package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
)

// TestAwsAccountExport verifies that an AWS account can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestAwsAccountExport() {
	testFramework := test.OctopusContainerTest{}

	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "3-awsaccount", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("..", "terraform", "3a-awsaccountds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := accounts.AccountsQuery{
		PartialName: "AWS Account",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Accounts.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an account called \"AWS Account\"")
	}
	resource := resources.Items[0].(*accounts.AmazonWebServicesAccount)

	if resource.AccessKey != "ABCDEFGHIJKLMNOPQRST" {
		t.Fatalf("The account must have an access key of \"ABCDEFGHIJKLMNOPQRST\"")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "/terraform", "3a-awsaccountds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
