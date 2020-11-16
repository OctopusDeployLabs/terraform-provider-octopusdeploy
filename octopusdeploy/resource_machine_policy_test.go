package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMachinePolicyImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_machine_policy." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccMachinePolicyCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testMachinePolicyBasic(localName, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMachinePolicyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccMachinePolicyCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testMachinePolicyExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testMachinePolicyBasic(localName, name),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr("data."+prefix, "name", name),
				),
				Config: testMachinePolicyDataSource(localName, name),
			},
		},
	})
}

func testMachinePolicyDataSource(localName string, name string) string {
	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "%s" {
		partial_name = "%s"
	}`, localName, name)
}

func testMachinePolicyBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name = "%s"
	}`, localName, name)
}

func testMachinePolicyExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		id := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.MachinePolicies.GetByID(id); err != nil {
			return err
		}

		return nil
	}
}

func testAccMachinePolicyCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		id := rs.Primary.ID
		machinePolicy, err := client.MachinePolicies.GetByID(id)
		if err == nil {
			if machinePolicy != nil {
				return fmt.Errorf("machine policy (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
