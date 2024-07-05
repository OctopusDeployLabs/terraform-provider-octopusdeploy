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

func TestAccOctopusDeployDockerContainerRegistry(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_docker_container_registry." + localName

	apiVersion := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedURI := "https://index.docker.io"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	registryPath := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDockerContainerRegistryCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testDockerContainerRegistryExists(prefix),
					resource.TestCheckResourceAttr(prefix, "api_version", apiVersion),
					resource.TestCheckResourceAttr(prefix, "feed_uri", feedURI),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "registry_path", registryPath),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testDockerContainerRegistryBasic(localName, apiVersion, feedURI, name, registryPath, username, password),
			},
		},
	})
}

func testDockerContainerRegistryBasic(localName string, apiVersion string, feedURI string, name string, registryPath string, username string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_docker_container_registry" "%s" {
		api_version = "%s"
		feed_uri = "%s"
		name = "%s"
		password = "%s"
		registry_path = "%s"
		username = "%s"
	}`, localName, apiVersion, feedURI, name, password, registryPath, username)
}

func testDockerContainerRegistryExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testDockerContainerRegistryCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_docker_container_registry" {
			continue
		}

		client := testAccProvider.Meta().(*client.Client)
		feed, err := client.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Docker Container Registry (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestDockerFeedResource verifies that a docker feed can be reimported with the correct settings
func TestDockerFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "11-dockerfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "11a-dockerfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := feeds.FeedsQuery{
			PartialName: "Docker",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Docker\"")
		}
		resource := resources.Items[0].(*feeds.DockerContainerRegistry)

		if resource.FeedType != "Docker" {
			t.Fatal("The feed must have a type of \"Docker\"")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.APIVersion != "v1" {
			t.Fatal("The feed must be have a API version of \"v1\"")
		}

		if resource.FeedURI != "https://index.docker.io" {
			t.Fatal("The feed must be have a feed uri of \"https://index.docker.io\"")
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
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "11a-dockerfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
