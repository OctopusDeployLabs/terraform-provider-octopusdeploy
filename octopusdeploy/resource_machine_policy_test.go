package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"strings"
)

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

// TestMachinePolicyResource verifies that a machine policies can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestMachinePolicyResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "27-machinepolicy", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := machines.MachinePoliciesQuery{
		PartialName: "Testing",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.MachinePolicies.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a machine policy called \"Testing\"")
	}
	resource := resources.Items[0]

	if resource.Description != "test machine policy" {
		t.Fatal("The machine policy must have a description of \"test machine policy\" (was \"" + resource.Description + "\")")
	}

	if resource.ConnectionConnectTimeout.Minutes() != 1 {
		t.Fatal("The machine policy must have a ConnectionConnectTimeout of \"00:01:00\" (was \"" + fmt.Sprint(resource.ConnectionConnectTimeout) + "\")")
	}

	if resource.ConnectionRetryCountLimit != 5 {
		t.Fatal("The machine policy must have a ConnectionRetryCountLimit of \"5\" (was \"" + fmt.Sprint(resource.ConnectionRetryCountLimit) + "\")")
	}

	if resource.ConnectionRetrySleepInterval.Seconds() != 1 {
		t.Fatal("The machine policy must have a ConnectionRetrySleepInterval of \"00:00:01\" (was \"" + fmt.Sprint(resource.ConnectionRetrySleepInterval) + "\")")
	}

	if resource.ConnectionRetryTimeLimit.Minutes() != 5 {
		t.Fatal("The machine policy must have a ConnectionRetryTimeLimit of \"00:05:00\" (was \"" + fmt.Sprint(resource.ConnectionRetryTimeLimit) + "\")")
	}

	if resource.MachineCleanupPolicy.DeleteMachinesElapsedTimeSpan.Minutes() != 20 {
		t.Fatal("The machine policy must have a DeleteMachinesElapsedTimeSpan of \"00:20:00\" (was \"" + fmt.Sprint(resource.MachineCleanupPolicy.DeleteMachinesElapsedTimeSpan) + "\")")
	}

	if resource.MachineCleanupPolicy.DeleteMachinesBehavior != "DeleteUnavailableMachines" {
		t.Fatal("The machine policy must have a MachineCleanupPolicy.DeleteMachinesBehavior of \"DeleteUnavailableMachines\" (was \"" + resource.MachineCleanupPolicy.DeleteMachinesBehavior + "\")")
	}

	if resource.MachineConnectivityPolicy.MachineConnectivityBehavior != "ExpectedToBeOnline" {
		t.Fatal("The machine policy must have a MachineConnectivityPolicy.MachineConnectivityBehavior of \"ExpectedToBeOnline\" (was \"" + resource.MachineConnectivityPolicy.MachineConnectivityBehavior + "\")")
	}

	if resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType != "Inline" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType of \"Inline\" (was \"" + resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType + "\")")
	}

	if *resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody != "" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody of \"\" (was \"" + *resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody + "\")")
	}

	if resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType != "Inline" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType of \"Inline\" (was \"" + resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType + "\")")
	}

	if strings.HasPrefix(*resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody, "$freeDiskSpaceThreshold") {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.ScriptBody to start with \"$freeDiskSpaceThreshold\" (was \"" + *resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.ScriptBody + "\")")
	}

	if resource.MachineHealthCheckPolicy.HealthCheckCronTimezone != "UTC" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckCronTimezone of \"UTC\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckCronTimezone + "\")")
	}

	if resource.MachineHealthCheckPolicy.HealthCheckCron != "" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckCron of \"\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckCron + "\")")
	}

	if resource.MachineHealthCheckPolicy.HealthCheckType != "RunScript" {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckType of \"RunScript\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckType + "\")")
	}

	if resource.MachineHealthCheckPolicy.HealthCheckInterval.Minutes() != 10 {
		t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckInterval of \"00:10:00\" (was \"" + fmt.Sprint(resource.MachineHealthCheckPolicy.HealthCheckInterval) + "\")")
	}

	if resource.MachineUpdatePolicy.CalamariUpdateBehavior != "UpdateAlways" {
		t.Fatal("The machine policy must have a MachineUpdatePolicy.CalamariUpdateBehavior of \"UpdateAlways\" (was \"" + resource.MachineUpdatePolicy.CalamariUpdateBehavior + "\")")
	}

	if resource.MachineUpdatePolicy.TentacleUpdateBehavior != "Update" {
		t.Fatal("The machine policy must have a MachineUpdatePolicy.TentacleUpdateBehavior of \"Update\" (was \"" + resource.MachineUpdatePolicy.CalamariUpdateBehavior + "\")")
	}

	if resource.MachineUpdatePolicy.KubernetesAgentUpdateBehavior != "NeverUpdate" {
		t.Fatal("The machine policy must have a MachineUpdatePolicy.KubernetesAgentUpdateBehavior of \"NeverUpdate\" (was \"" + resource.MachineUpdatePolicy.CalamariUpdateBehavior + "\")")
	}
}
