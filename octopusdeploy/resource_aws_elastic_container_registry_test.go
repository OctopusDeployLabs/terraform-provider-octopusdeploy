package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"os"
	"path/filepath"
	"testing"
)

// TestEcrFeedResource verifies that a ecr feed can be reimported with the correct settings
func TestEcrFeedResource(t *testing.T) {
	if os.Getenv("ECR_ACCESS_KEY") == "" {
		t.Fatal("The ECR_ACCESS_KEY environment variable must be set a valid AWS access key")
	}

	if os.Getenv("ECR_SECRET_KEY") == "" {
		t.Fatal("The ECR_SECRET_KEY environment variable must be set a valid AWS secret key")
	}

	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "12-ecrfeed", []string{
		"-var=feed_ecr_access_key=" + os.Getenv("ECR_ACCESS_KEY"),
		"-var=feed_ecr_secret_key=" + os.Getenv("ECR_SECRET_KEY"),
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "12a-ecrfeedds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := feeds.FeedsQuery{
		PartialName: "ECR",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Feeds.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an feed called \"ECR\"")
	}
	resource := resources.Items[0].(*feeds.AwsElasticContainerRegistry)

	if resource.FeedType != "AwsElasticContainerRegistry" {
		t.Fatal("The feed must have a type of \"AwsElasticContainerRegistry\" (was \"" + resource.FeedType + "\"")
	}

	if resource.AccessKey != os.Getenv("ECR_ACCESS_KEY") {
		t.Fatal("The feed must have a access key of \"" + os.Getenv("ECR_ACCESS_KEY") + "\" (was \"" + resource.AccessKey + "\"")
	}

	if resource.Region != "us-east-1" {
		t.Fatal("The feed must have a region of \"us-east-1\" (was \"" + resource.Region + "\"")
	}

	foundExecutionTarget := false
	foundNotAcquired := false
	for _, o := range resource.PackageAcquisitionLocationOptions {
		if o == "ExecutionTarget" {
			foundExecutionTarget = true
		}

		if o == "NotAcquired" {
			foundNotAcquired = true
		}
	}

	if !(foundExecutionTarget && foundNotAcquired) {
		t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"NotAcquired\"")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "12a-ecrfeedds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
