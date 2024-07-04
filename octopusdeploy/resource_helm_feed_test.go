package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		CheckDestroy: testHelmFeedCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		client := testAccProvider.Meta().(*client.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
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

		client := testAccProvider.Meta().(*client.Client)
		feed, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Helm feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestHelmFeedResource verifies that a helm feed can be reimported with the correct settings
func TestHelmFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "10-helmfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "10a-helmfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := feeds.FeedsQuery{
			PartialName: "Helm",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
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
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "10a-helmfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
