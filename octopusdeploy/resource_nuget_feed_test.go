package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployNuGetFeedBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_nuget_feed." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedURI := "http://test.com"
	isEnhancedMode := true
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployNuGetFeedDestroy,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployNuGetFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "is_enhanced_mode", strconv.FormatBool(isEnhancedMode)),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testAccNuGetFeed(localName, name, feedURI, username, password, isEnhancedMode),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployNuGetFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "is_enhanced_mode", strconv.FormatBool(isEnhancedMode)),
					resource.TestCheckResourceAttr(prefix, "name", updatedName),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testAccNuGetFeed(localName, updatedName, feedURI, username, password, isEnhancedMode),
			},
		},
	})
}

func testAccNuGetFeed(localName string, name string, feedURI string, username string, password string, isEnhancedMode bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_nuget_feed" "%s" {
		feed_uri         = "%s"
		is_enhanced_mode = %v
		name             = "%s"
		password         = "%s"
		username         = "%s"
	}`, localName, feedURI, isEnhancedMode, name, password, username)
}

func testOctopusDeployNuGetFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testOctopusDeployNuGetFeedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_nuget_feed" {
			continue
		}

		_, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestNugetFeedResource verifies that a nuget feed can be reimported with the correct settings
func TestNugetFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "14-nugetfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "14a-nugetfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := feeds.FeedsQuery{
			PartialName: "Nuget",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Nuget\"")
		}
		resource := resources.Items[0].(*feeds.NuGetFeed)

		if resource.FeedType != "NuGet" {
			t.Fatal("The feed must have a type of \"NuGet\"")
		}

		if !resource.EnhancedMode {
			t.Fatal("The feed must have enhanced mode set to true")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.DownloadAttempts != 5 {
			t.Fatal("The feed must be have a downloads attempts set to \"5\"")
		}

		if resource.DownloadRetryBackoffSeconds != 10 {
			t.Fatal("The feed must be have a downloads retry backoff set to \"10\"")
		}

		if resource.FeedURI != "https://index.docker.io" {
			t.Fatal("The feed must be have a feed uri of \"https://index.docker.io\"")
		}

		foundExecutionTarget := false
		foundServer := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "Server" {
				foundServer = true
			}
		}

		if !(foundExecutionTarget && foundServer) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"Server\"")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "14a-nugetfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
