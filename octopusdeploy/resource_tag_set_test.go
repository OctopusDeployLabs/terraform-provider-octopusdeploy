package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployTagSetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := constOctopusDeployTagSet + "." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagColor := "#6e6e6e"
	tagDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSortOrder := acctest.RandIntRange(0, 10)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testTagSetDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testTagSetMinimal(localName, name),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tag.0.color", tagColor),
					resource.TestCheckResourceAttr(prefix, "tag.0.description", tagDescription),
					resource.TestCheckResourceAttr(prefix, "tag.0.name", tagName),
					resource.TestCheckResourceAttr(prefix, "tag.0.sort_order", strconv.Itoa(tagSortOrder)),
				),
				Config: testTagSetComplete(localName, name, tagColor, tagDescription, tagName, tagSortOrder),
			},
		},
	})
}

func testTagSetMinimal(localName string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
	}`, constOctopusDeployTagSet, localName, name)
}

func testTagSetComplete(localName string, name string, tagColor string, tagDescription string, tagName string, tagSortOrder int) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
		tag {
			color = "%s"
			description = "%s"
			name = "%s"
			sort_order = "%v"
		}
	}`, constOctopusDeployTagSet, localName, name, tagColor, tagDescription, tagName, tagSortOrder)
}

func testTagSetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		tagSetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.TagSets.GetByID(tagSetID); err != nil {
			return err
		}

		return nil
	}
}

func testTagSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		tagSetID := rs.Primary.ID
		tagSet, err := client.TagSets.GetByID(tagSetID)
		if err == nil {
			if tagSet != nil {
				return fmt.Errorf("tag set (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
