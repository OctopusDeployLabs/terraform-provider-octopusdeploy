package octopusdeploy

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccDataSourceWorkerPools() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_worker_pools.%s", localName)
	partialName := "W"

	resource.Test(suite.T(), resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkerPoolsConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkerPoolsDataSourceID(name),
					resource.TestCheckResourceAttrSet(name, "worker_pools.#"),
				)},
		},
	})
}

func testAccCheckWorkerPoolsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find worker pools data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("worker pools ID not set")
		}
		return nil
	}
}

func testAccDataSourceWorkerPoolsConfig(localName string, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_worker_pools" "%s" {
		partial_name = "%s"
	}`, localName, partialName)
}
