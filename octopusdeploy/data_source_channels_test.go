package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceChannels(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_channels.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceChannelsConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChannelsDataSourceID(name),
				)},
		},
	})
}

func testAccCheckChannelsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Channels data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Channels source ID not set")
		}
		return nil
	}
}

func testAccDataSourceChannelsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_channels" "%s" {
		take = %v
	}`, localName, take)
}
