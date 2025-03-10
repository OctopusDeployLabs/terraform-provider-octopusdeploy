package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
)

func TestTagSetAndTag(t *testing.T) {
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetPrefix := "octopusdeploy_tag_set." + tagSetName
	tagSetDescription := "TagSet Description" + tagSetName

	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagPrefix := "octopusdeploy_tag." + tagName
	tagColor := "#6e6e6e"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testTagSetConfig(tagSetName, tagSetDescription),
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(tagSetPrefix),
					resource.TestCheckResourceAttr(tagSetPrefix, "name", tagSetName),
					resource.TestCheckResourceAttr(tagSetPrefix, "description", tagSetDescription),
				),
			},
			{
				Config: testTagSetAndTagConfig(tagSetName, tagSetDescription, tagName, tagColor),
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(tagSetPrefix),
					testTagExists(tagPrefix),
					resource.TestCheckResourceAttr(tagSetPrefix, "name", tagSetName),
					resource.TestCheckResourceAttr(tagSetPrefix, "description", tagSetDescription),
					resource.TestCheckResourceAttr(tagPrefix, "name", tagName),
					resource.TestCheckResourceAttr(tagPrefix, "color", tagColor),
					resource.TestCheckResourceAttrSet(tagPrefix, "id"),
					resource.TestCheckResourceAttrSet(tagPrefix, "tag_set_space_id"),
					resource.TestCheckResourceAttrSet(tagPrefix, "tag_set_id"),
					resource.TestCheckResourceAttrPair(tagPrefix, "tag_set_id", tagSetPrefix, "id"),
				),
			},
		},
	})
}

func testTagSetConfig(name, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "%s" {
		  name        = "%s"
		  description = "%s"
		}`, name, name, description)
}

func testTagSetAndTagConfig(tagSetName, tagSetDescription, tagName, tagColor string) string {
	var tfConfig = fmt.Sprintf(`
    resource "octopusdeploy_tag_set" "%s" {
      name        = "%s"
      description = "%s"
    }
    
    resource "octopusdeploy_tag" "%s" {
      name        = "%s"
      color       = "%s"
      description = "Test tag"
      tag_set_id  = octopusdeploy_tag_set.%s.id
    }`, tagSetName, tagSetName, tagSetDescription, tagName, tagName, tagColor, tagSetName)
	return tfConfig
}

func testTagSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if _, err := tagsets.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err != nil {
			return fmt.Errorf("error retrieving tag set: %s", err)
		}

		return nil
	}
}

func testTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		tagSetID := rs.Primary.Attributes["tag_set_id"]
		spaceID := rs.Primary.Attributes["tag_set_space_id"]
		tagSet, err := tagsets.GetByID(octoClient, spaceID, tagSetID)
		if err != nil {
			return fmt.Errorf("error retrieving tag set: %s", err)
		}

		for _, tag := range tagSet.Tags {
			if tag.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("tag not found in tag set")
	}
}

func TestExpandTagSet(t *testing.T) {
	name := "Test Tag Set"
	description := "This is a test tag set"
	sortOrder := int64(10)
	spaceID := "Spaces-1"

	tagSetModel := schemas.TagSetResourceModel{
		Name:        types.StringValue(name),
		Description: types.StringValue(description),
		SortOrder:   types.Int64Value(sortOrder),
		SpaceID:     types.StringValue(spaceID),
	}

	tagSet := expandTagSet(tagSetModel)

	require.Equal(t, name, tagSet.Name)
	require.Equal(t, description, tagSet.Description)
	require.Equal(t, int32(sortOrder), tagSet.SortOrder)
	require.Equal(t, spaceID, tagSet.SpaceID)
}
