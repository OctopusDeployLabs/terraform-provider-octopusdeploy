package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
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
		CheckDestroy:             testDockerContainerRegistryCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
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

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Docker Container Registry (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
