package octopusdeploy

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
// 	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// func TestAccMachinePolicyImportBasic(t *testing.T) {
// 	options := test.NewTestOptions()
// 	resourceName := "octopusdeploy_machine_policy." + options.LocalName

// 	resource.Test(t, resource.TestCase{
// 		CheckDestroy: testAccMachinePolicyCheckDestroy,
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testMachinePolicyBasic(options),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// func TestAccMachinePolicyIssue230(t *testing.T) {
// 	options := test.NewTestOptions()

// 	resource.Test(t, resource.TestCase{
// 		CheckDestroy: testAccMachinePolicyCheckDestroy,
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testMachinePolicyIssue230(options),
// 			},
// 		},
// 	})
// }

// func TestAccMachinePolicyBasic(t *testing.T) {
// 	options := test.NewTestOptions()
// 	resourceName := "octopusdeploy_machine_policy." + options.LocalName

// 	resource.Test(t, resource.TestCase{
// 		CheckDestroy: testAccMachinePolicyCheckDestroy,
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Check: resource.ComposeTestCheckFunc(
// 					testMachinePolicyExists(resourceName),
// 					resource.TestCheckResourceAttrSet(resourceName, "id"),
// 					resource.TestCheckResourceAttr(resourceName, "name", options.Name),
// 				),
// 				Config: testMachinePolicyBasic(options),
// 			},
// 			{
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet("data.octopusdeploy_machine_policies."+options.LocalName, "id"),
// 					resource.TestCheckResourceAttr("data.octopusdeploy_machine_policies."+options.LocalName, "machine_policies.#", "1"),
// 				),
// 				Config: testMachinePolicyDataSource(options),
// 			},
// 		},
// 	})
// }

// func testMachinePolicyDataSource(options *test.TestOptions) string {
// 	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "%s" {
// 		partial_name = "%s"
// 	}`, options.LocalName, options.Name)
// }

// func testMachinePolicyBasic(options *test.TestOptions) string {
// 	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
// 		name = "%s"
// 	}`, options.LocalName, options.Name)
// }

// func testMachinePolicyIssue230(options *test.TestOptions) string {
// 	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
// 		name = "%s"

// 		machine_connectivity_policy {
// 		  machine_connectivity_behavior = "ExpectedToBeOnline"
// 		}

// 		machine_cleanup_policy {
// 		  delete_machines_behavior         = "DeleteUnavailableMachines"
// 		  delete_machines_elapsed_timespan = 3600000000000
// 		}

// 		machine_health_check_policy {
// 		  health_check_type     = "OnlyConnectivity"
// 		  health_check_interval = 15000000000

// 		  bash_health_check_policy {
// 			run_type = "InheritFromDefault"
// 		  }

// 		  powershell_health_check_policy {
// 			run_type = "InheritFromDefault"
// 		}
// 	  }

// 		machine_update_policy {
// 		  calamari_update_behavior = "UpdateAlways"
// 		  tentacle_update_behavior = "NeverUpdate"
// 		}
// 	}`, options.LocalName, options.Name)
// }

// func testMachinePolicyExists(prefix string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		client := testAccProvider.Meta().(*client.Client)
// 		id := s.RootModule().Resources[prefix].Primary.ID
// 		if _, err := client.MachinePolicies.GetByID(id); err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// }

// func testAccMachinePolicyCheckDestroy(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*client.Client)
// 	for _, rs := range s.RootModule().Resources {
// 		id := rs.Primary.ID
// 		machinePolicy, err := client.MachinePolicies.GetByID(id)
// 		if err == nil {
// 			if machinePolicy != nil {
// 				return fmt.Errorf("machine policy (%s) still exists", rs.Primary.ID)
// 			}
// 		}
// 	}

// 	return nil
// }
