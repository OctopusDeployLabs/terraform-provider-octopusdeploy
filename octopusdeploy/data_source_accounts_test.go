package octopusdeploy

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccDataSourceAccounts() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_accounts.%s", localName)
	take := 10

	resource.Test(suite.T(), resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountsDataSourceID(name),
				)},
		},
	})
}

func testAccCheckAccountsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Accounts data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Accounts source ID not set")
		}
		return nil
	}
}

func testAccDataSourceAccountsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_accounts" "%s" {
		take = %v
	}`, localName, take)
}
