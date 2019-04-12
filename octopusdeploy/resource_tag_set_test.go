package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployTagSetBasic(t *testing.T) {
	const tagSetPrefix = "octopusdeploy_tag_set.foo"
	const tagSetName = "Testing one two three"
	const tagName1 = "tagName1"
	const tagName2 = "tagName2"
	const tagColor1 = "#6e6e6e"
	const tagColor2 = "#6e6e6f"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployTagSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTagSettBasic(tagSetName, tagName1, tagColor1, tagName2, tagColor2),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployTagSetExists(tagSetPrefix),
					resource.TestCheckResourceAttr(
						tagSetPrefix, "name", tagSetName),
					resource.TestCheckResourceAttr(
						tagSetPrefix, "tag.0.name", tagName1),
					resource.TestCheckResourceAttr(
						tagSetPrefix, "tag.0.color", tagColor1),
					resource.TestCheckResourceAttr(
						tagSetPrefix, "tag.1.name", tagName2),
					resource.TestCheckResourceAttr(
						tagSetPrefix, "tag.1.color", tagColor2),
				),
			},
		},
	})
}

func testTagSettBasic(name, tagName1 string, tagColor1 string, tagName2 string, tagColor2 string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_tag_set" "foo" {
			name = "%s"

			tag {
				name = "%s"
				color = "%s"
			}

			tag {
				name = "%s"
				color = "%s"
			}
		}
		`,
		name, tagName1, tagColor1, tagName2, tagColor2,
	)
}

func testOctopusDeployTagSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existstagSetHelper(s, client)
	}
}

func existstagSetHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.TagSet.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving tagSet %s", err)
		}
	}
	return nil
}

func testOctopusDeployTagSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroytagSetHelper(s, client)
}

func destroytagSetHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.TagSet.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving tagSet %s", err)
		}
		return fmt.Errorf("TagSet still exists")
	}
	return nil
}
