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
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
)

func TestAccOctopusDeployTagSetAndTag(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetPrefix := "octopusdeploy_tag_set." + localName
	tagPrefix := "octopusdeploy_tag." + localName

	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagColor := "#6e6e6e"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testTagSetConfig(localName, tagSetName, tagSetDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccOctopusDeployTagSetExists(tagSetPrefix),
					resource.TestCheckResourceAttr(tagSetPrefix, "name", tagSetName),
					resource.TestCheckResourceAttr(tagSetPrefix, "description", tagSetDescription),
				),
			},
			{
				Config: testTagSetAndTagConfig(localName, tagSetName, tagSetDescription, tagName, tagColor),
				Check: resource.ComposeTestCheckFunc(
					testAccOctopusDeployTagSetExists(tagSetPrefix),
					testAccOctopusDeployTagExists(tagPrefix),
					resource.TestCheckResourceAttr(tagSetPrefix, "name", tagSetName),
					resource.TestCheckResourceAttr(tagSetPrefix, "description", tagSetDescription),
					resource.TestCheckResourceAttr(tagPrefix, "name", tagName),
					resource.TestCheckResourceAttr(tagPrefix, "color", tagColor),
					resource.TestCheckResourceAttrSet(tagPrefix, "id"),
					resource.TestCheckResourceAttrSet(tagPrefix, "space_id"),
					resource.TestCheckResourceAttrSet(tagPrefix, "tag_set_id"),
					resource.TestCheckResourceAttrPair(tagPrefix, "tag_set_id", tagSetPrefix, "id"),
				),
			},
		},
	})
}

func testTagSetConfig(localName, name, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "%s" {
		  name        = "%s"
		  description = "%s"
		}`, localName, name, description)
}

func testTagSetAndTagConfig(localName, tagSetName, tagSetDescription, tagName, tagColor string) string {
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

func testAccOctopusDeployTagSetExists(n string) resource.TestCheckFunc {
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

func testAccOctopusDeployTagExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		tagSetID := rs.Primary.Attributes["tag_set_id"]
		tagSet, err := tagsets.GetByID(octoClient, rs.Primary.Attributes["space_id"], tagSetID)
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

// TestTagSetResource verifies that a tag set can be reimported with the correct settings
func TestTagSetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "21-tagset", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := tagsets.TagSetsQuery{
		PartialName: "tag1",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.TagSets.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a tag set called \"tag1\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test tagset" {
		t.Fatal("The tag set must be have a description of \"Test tagset\" (was \"" + resource.Description + "\")")
	}

	if resource.SortOrder != 0 {
		t.Fatal("The tag set must be have a sort order of \"0\" (was \"" + fmt.Sprint(resource.SortOrder) + "\")")
	}

	tagAFound := false
	for _, u := range resource.Tags {
		if u.Name == "a" {
			tagAFound = true

			if u.Description != "tag a" {
				t.Fatal("The tag a must be have a description of \"tag a\" (was \"" + u.Description + "\")")
			}

			if u.Color != "#333333" {
				t.Fatal("The tag a must be have a color of \"#333333\" (was \"" + u.Color + "\")")
			}

			if u.SortOrder != 2 {
				t.Fatal("The tag a must be have a sort order of \"2\" (was \"" + fmt.Sprint(u.SortOrder) + "\")")
			}
		}
	}

	if !tagAFound {
		t.Fatal("Tag Set must have an tag called \"a\"")
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
