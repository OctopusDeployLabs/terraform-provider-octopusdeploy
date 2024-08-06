package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployTagSetBasic(t *testing.T) {
	internalTest.SkipCI(t, " Error: Unsupported block type")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_tag_set." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagColor := "#6e6e6e"
	tagDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testTagSetDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tag.#", "0"),
				),
				Config: testTagSetMinimal(localName, name),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "tag.0.id"),
					resource.TestCheckResourceAttr(prefix, "tag.0.color", tagColor),
					resource.TestCheckResourceAttr(prefix, "tag.0.description", tagDescription),
					resource.TestCheckResourceAttr(prefix, "tag.0.name", tagName),
				),
				Config: testTagSetComplete(localName, name, tagColor, tagDescription, tagName),
			},
		},
	})
}

func testTagSetMinimal(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tag_set" "%s" {
		name = "%s"
	}`, localName, name)
}

func testTagSetComplete(localName string, name string, tagColor string, tagDescription string, tagName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tag_set" "%s" {
		name = "%s"
		tag {
			color = "%s"
			description = "%s"
			name = "%s"
		}
	}`, localName, name, tagColor, tagDescription, tagName)
}

func testTagSetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tagSetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.TagSets.GetByID(tagSetID); err != nil {
			return err
		}

		return nil
	}
}

func testTagSetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		tagSetID := rs.Primary.ID
		tagSet, err := octoClient.TagSets.GetByID(tagSetID)
		if err == nil {
			if tagSet != nil {
				return fmt.Errorf("tag set (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

// TestTagSetResource verifies that a tag set can be reimported with the correct settings
func TestTagSetResource(t *testing.T) {
	internalTest.SkipCI(t, " Error: Unsupported block type")
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
