package octopusdeploy

import (
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"net/url"
	"path/filepath"
)

func (suite *IntegrationTestSuite) TestPollingSubscriptionIdResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	_, err := testFramework.Act(t, octoContainer, "../terraform", "56-pollingsubscriptionid", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	baseIdLookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "56-pollingsubscriptionid"), "base_id")
	if err != nil {
		t.Fatal(err.Error())
	}

	basePollingUriLookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "56-pollingsubscriptionid"), "base_polling_uri")
	if err != nil {
		t.Fatal(err.Error())
	}

	parsedUri, err := url.Parse(basePollingUriLookup)
	if parsedUri.Scheme != "poll" {
		t.Fatalf("The polling URI scheme must be \"poll\" but instead received %s", parsedUri.Scheme)
	}

	if parsedUri.Host != baseIdLookup {
		t.Fatalf("The polling URI host must be the Subscription ID but instead received %s", parsedUri.Host)
	}
}
