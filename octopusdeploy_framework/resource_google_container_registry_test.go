package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

type googleFeedTestData struct {
	name         string
	uri          string
	registryPath string
	apiVersion   string
	username     string
	password     string
}

func TestAccOctopusDeployGoogleFeed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_google_container_registry." + localName
	createData := googleFeedTestData{
		name:         acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		uri:          "https://cloud.artifact.google.test",
		registryPath: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
		apiVersion:   acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
		username:     acctest.RandStringFromCharSet(16, acctest.CharSetAlpha),
		password:     acctest.RandStringFromCharSet(300, acctest.CharSetAlpha),
	}
	updateData := googleFeedTestData{
		name:         createData.name + "-updated",
		uri:          "https://testcloud.artifact.google.updated",
		registryPath: createData.registryPath + "-updated",
		apiVersion:   createData.apiVersion + "-updated",
		username:     createData.username + "-updated",
		password:     createData.password + "-updated",
	}
	withMinimumData := googleFeedTestData{
		name: "Google Registry Minimum",
		uri:  "https://testcloud.artifact.google.minimum",
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testGoogleFeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testGoogleFeedBasic(createData, localName),
				Check:  testAssertGoogleFeedAttributes(createData, prefix),
			},
			{
				Config: testGoogleFeedBasic(updateData, localName),
				Check:  testAssertGoogleFeedAttributes(updateData, prefix),
			},
			{
				Config: testGoogleFeedWithMinimumData(withMinimumData, localName),
				Check:  testAssertGoogleFeedMinimumAttributes(withMinimumData, prefix),
			},
		},
	})
}

func testGoogleFeedBasic(data googleFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_google_container_registry" "%s" {
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

func testGoogleFeedWithMinimumData(data googleFeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_google_container_registry" "%s" {
			name			= "%s"
			feed_uri		= "%s"
		}
	`,
		localName,
		data.name,
		data.uri,
	)
}

func testAssertGoogleFeedAttributes(expected googleFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "feed_uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "registry_path", expected.registryPath),
		resource.TestCheckResourceAttr(prefix, "api_version", expected.apiVersion),
		resource.TestCheckResourceAttr(prefix, "username", expected.username),
		resource.TestCheckResourceAttr(prefix, "password", expected.password),
	)
}

func testAssertGoogleFeedMinimumAttributes(expected googleFeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "feed_uri", expected.uri),
		resource.TestCheckNoResourceAttr(prefix, "registry_path"),
		resource.TestCheckNoResourceAttr(prefix, "api_version"),
		resource.TestCheckNoResourceAttr(prefix, "username"),
		resource.TestCheckNoResourceAttr(prefix, "password"),
	)
}

func testGoogleFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_google_container_registry_feed" {
			continue
		}

		feed, err := feeds.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("google container registry feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
