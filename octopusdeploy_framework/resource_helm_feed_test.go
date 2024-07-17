package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"path/filepath"
	"testing"
)

func TestAccOctopusDeployHelmFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_helm_feed." + localName

	feedURI := "https://charts.helm.sh/stable"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testHelmFeedCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testHelmFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testHelmFeedBasic(localName, feedURI, name, username, password),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testHelmFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", updatedName),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testHelmFeedBasic(localName, feedURI, updatedName, username, password),
			},
		},
	})
}

func testHelmFeedBasic(localName string, feedURI string, name string, username string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_helm_feed" "%s" {
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		username = "%s"
	}`, localName, feedURI, name, password, username)
}

func testHelmFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testHelmFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_helm_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Helm feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestHelmFeedResource verifies that a helm feed can be reimported with the correct settings
func TestHelmFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "10-helmfeed", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "10a-helmfeedds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := feeds.FeedsQuery{
		PartialName: "Helm",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Feeds.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an feed called \"Helm\"")
	}
	resource := resources.Items[0].(*feeds.HelmFeed)

	if resource.FeedType != "Helm" {
		t.Fatal("The feed must have a type of \"Helm\"")
	}

	if resource.Username != "username" {
		t.Fatal("The feed must have a username of \"username\"")
	}

	if resource.FeedURI != "https://charts.helm.sh/stable/" {
		t.Fatal("The feed must be have a URI of \"https://charts.helm.sh/stable/\"")
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
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "10a-helmfeedds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
