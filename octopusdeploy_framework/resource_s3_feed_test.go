package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

type s3FeedTestData struct {
	name                  string
	useMachineCredentials bool
	accessKey             string
	secretKey             string
}

func TestAccOctopusDeployS3Feed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_s3_feed." + localName
	createData := s3FeedTestData{
		name:                  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		useMachineCredentials: false,
		accessKey:             acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		secretKey:             acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum),
	}
	updateData := s3FeedTestData{
		name:                  createData.name + "-updated",
		useMachineCredentials: true,
		accessKey:             createData.accessKey + "-changed",
		secretKey:             createData.secretKey + "-generated",
	}
	withoutKeysData := s3FeedTestData{
		name:                  "AWS S3 Without Keys",
		useMachineCredentials: true,
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testS3FeedCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testS3FeedBasic(createData, localName),
				Check:  testAssertS3FeedAttributes(createData, prefix),
			},
			{
				Config: testS3FeedBasic(updateData, localName),
				Check:  testAssertS3FeedAttributes(updateData, prefix),
			},
			{
				Config: testS3FeedWithoutKeys(withoutKeysData, localName),
				Check:  testAssertS3FeedAttributes(withoutKeysData, prefix),
			},
		},
	})
}

func testAssertS3FeedAttributes(expected s3FeedTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "use_machine_credentials", strconv.FormatBool(expected.useMachineCredentials)),
		resource.TestCheckResourceAttr(prefix, "access_key", expected.accessKey),
		resource.TestCheckResourceAttr(prefix, "secret_key", expected.secretKey),
	)
}

func testS3FeedBasic(data s3FeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_s3_feed" "%s" {
			name		= "%s"
			use_machine_credentials	= "%s"
			access_key	= "%s"
			secret_key	= "%s"
		}
	`,
		localName,
		data.name,
		strconv.FormatBool(data.useMachineCredentials),
		data.accessKey,
		data.secretKey,
	)
}

func testS3FeedWithoutKeys(data s3FeedTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_s3_feed" "%s" {
			name		= "%s"
			use_machine_credentials	= "%s"
		}
	`,
		localName,
		data.name,
		strconv.FormatBool(data.useMachineCredentials),
	)
}

func testS3FeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_s3_feed" {
			continue
		}

		feed, err := feeds.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("s3 feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
