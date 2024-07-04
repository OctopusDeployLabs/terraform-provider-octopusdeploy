package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"net/url"
	"path/filepath"
	"testing"
)

func TestPollingSubscriptionIdResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		_, err := testFramework.Act(t, container, "../terraform", "56-pollingsubscriptionid", []string{})

		if err != nil {
			return err
		}

		baseIdLookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "56-pollingsubscriptionid"), "base_id")
		if err != nil {
			return err
		}

		basePollingUriLookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "56-pollingsubscriptionid"), "base_polling_uri")
		if err != nil {
			return err
		}

		parsedUri, err := url.Parse(basePollingUriLookup)
		if parsedUri.Scheme != "poll" {
			t.Fatalf("The polling URI scheme must be \"poll\" but instead received %s", parsedUri.Scheme)
		}

		if parsedUri.Host != baseIdLookup {
			t.Fatalf("The polling URI host must be the Subscription ID but instead received %s", parsedUri.Host)
		}

		return nil
	})
}
