package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

type azureFeedTestData struct {
	name         string
	uri          string
	registryPath string
	apiVersion   string
	username     string
	password     string
}

func TestAccOctopusDeployAzureFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_container_registry." + localName
	createData := azureFeedTestData{
		name:         acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		uri:          "https://azure.io.test",
		registryPath: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
		apiVersion:   acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
		username:     acctest.RandStringFromCharSet(16, acctest.CharSetAlpha),
		password:     acctest.RandStringFromCharSet(300, acctest.CharSetAlpha),
	}
	updateData := azureFeedTestData{
		name:         createData.name + "-updated",
		uri:          "https://azure.io.test.updated",
		registryPath: createData.registryPath + "-updated",
		apiVersion:   createData.apiVersion + "-updated",
		username:     createData.username + "-updated",
		password:     createData.password + "-updated",
	}
	withMinimumData := azureFeedTestData{
		name: "Azure Registry Minimum",
		uri:  "https://test-azure.minimum",
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testAzureFeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAzureFeedBasic(createData, localName),
				Check:  testAssertAzureFeedAttributes(createData, prefix),
			},
			{
				Config: testAzureFeedBasic(updateData, localName),
				Check:  testAssertAzureFeedAttributes(updateData, prefix),
			},
			{
				Config: testAzureFeedWithMinimumData(withMinimumData, localName),
				Check:  testAssertAzureFeedMinimumAttributes(withMinimumData, prefix),
			},
		},
	})
}

func testAzureFeedBasic(data azureFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_azure_container_registry" "%s" {
			name			= "%s"
			feed_uri		= "%s"
			registry_path	= "%s"
			api_version		= "%s"
			username		= "%s"
			password		= "%s"
		}
	`,
		localName,
		data.name,
		data.uri,
		data.registryPath,
		data.apiVersion,
		data.username,
		data.password,
	)
}

func testAzureFeedWithMinimumData(data azureFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_azure_container_registry" "%s" {
			name			= "%s"
			feed_uri		= "%s"
		}
	`,
		localName,
		data.name,
		data.uri,
	)
}

func testAssertAzureFeedAttributes(expected azureFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "feed_uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "registry_path", expected.registryPath),
		resource.TestCheckResourceAttr(prefix, "api_version", expected.apiVersion),
		resource.TestCheckResourceAttr(prefix, "username", expected.username),
		resource.TestCheckResourceAttr(prefix, "password", expected.password),
	)
}

func testAssertAzureFeedMinimumAttributes(expected azureFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "feed_uri", expected.uri),
		resource.TestCheckNoResourceAttr(prefix, "registry_path"),
		resource.TestCheckNoResourceAttr(prefix, "api_version"),
		resource.TestCheckNoResourceAttr(prefix, "username"),
		resource.TestCheckNoResourceAttr(prefix, "password"),
	)
}

func testAzureFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_container_registry_feed" {
			continue
		}

		feed, err := feeds.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("azure container registry feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
