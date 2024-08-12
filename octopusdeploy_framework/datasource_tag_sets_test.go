package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccDataSourceTagSets(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetResourceName := fmt.Sprintf("octopusdeploy_tag_set.%s", localName)
	dataSourceName := fmt.Sprintf("data.octopusdeploy_tag_sets.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Create a tag set
			{
				Config: testAccTagSetConfig(localName, tagSetName),
				Check: resource.ComposeTestCheckFunc(
					testTagSetExists(tagSetResourceName),
					resource.TestCheckResourceAttr(tagSetResourceName, "name", tagSetName),
				),
			},
			// Query the created tag set using the data source
			{
				Config: testAccTagSetConfig(localName, tagSetName) + testAccDataSourceTagSetsConfig(localName, tagSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagSetsDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.id", tagSetResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag_sets.0.name", tagSetResourceName, "name"),
				),
			},
		},
	})
}

func testAccCheckTagSetsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cannot find TagSets data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("TagSets data source ID not set")
		}
		return nil
	}
}

func testAccTagSetConfig(localName, tagSetName string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_tag_set" "%s" {
    name        = "%s"
    description = "Test tag set"
}
`, localName, tagSetName)
}

func testAccDataSourceTagSetsConfig(localName, tagSetName string) string {
	return fmt.Sprintf(`
data "octopusdeploy_tag_sets" "%s" {
    partial_name = "%s"
    skip         = 0
    take         = 10
}
`, localName, tagSetName)
}
