package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTag(t *testing.T) {
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagColor := "#6e6e6e"
	tagResourceName := "octopusdeploy_tag." + tagName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTagConfig(tagSetName, tagName, tagColor),
				Check: resource.ComposeTestCheckFunc(
					testTagExists(tagResourceName),
					resource.TestCheckResourceAttr(tagResourceName, "name", tagName),
					resource.TestCheckResourceAttr(tagResourceName, "color", tagColor),
					resource.TestCheckResourceAttrSet(tagResourceName, "id"),
					resource.TestCheckResourceAttrSet(tagResourceName, "tag_set_id"),
					resource.TestCheckResourceAttrSet(tagResourceName, "tag_set_space_id"),
				),
			},
			{
				Config: testTagConfigUpdate(tagSetName, tagName, "#ff0000"),
				Check: resource.ComposeTestCheckFunc(
					testTagExists(tagResourceName),
					resource.TestCheckResourceAttr(tagResourceName, "name", tagName),
					resource.TestCheckResourceAttr(tagResourceName, "color", "#ff0000"),
					resource.TestCheckResourceAttrSet(tagResourceName, "id"),
					resource.TestCheckResourceAttrSet(tagResourceName, "tag_set_id"),
					resource.TestCheckResourceAttrSet(tagResourceName, "tag_set_space_id"),
				),
			},
			{
				ResourceName:      tagResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTagImportStateIdFunc(tagResourceName),
			},
		},
	})
}

func testTagConfig(tagSetName, tagName, tagColor string) string {
	var tfConfig = fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "%s" {
			name        = "%s"
			description = "Test tag set"
		}

		resource "octopusdeploy_tag" "%s" {
			name        = "%s"
			color       = "%s"
			description = "Test tag"
			tag_set_id  = octopusdeploy_tag_set.%s.id
		}
	`, tagSetName, tagSetName, tagName, tagName, tagColor, tagSetName)
	return tfConfig
}

func testTagConfigUpdate(tagSetName, tagName, tagColor string) string {
	var tfConfig = fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "%s" {
			name        = "%s"
			description = "Test tag set"
		}

		resource "octopusdeploy_tag" "%s" {
			name        = "%s"
			color       = "%s"
			description = "Updated test tag"
			tag_set_id  = octopusdeploy_tag_set.%s.id
		}
	`, tagSetName, tagSetName, tagName, tagName, tagColor, tagSetName)
	return tfConfig
}

func testAccTagDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tag" {
			continue
		}

		tagSetID := rs.Primary.Attributes["tag_set_id"]
		tagSet, err := tagsets.GetByID(octoClient, rs.Primary.Attributes["space_id"], tagSetID)
		if err != nil {
			return nil // If the tag set is gone, the tag is gone too
		}

		for _, tag := range tagSet.Tags {
			if tag.ID == rs.Primary.ID {
				return fmt.Errorf("Tag still exists")
			}
		}
	}

	return nil
}

func testAccTagImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		tagID := rs.Primary.ID

		return tagID, nil
	}
}
