package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployMachineBasic(t *testing.T) {
	const tfVarPrefix = "octopusdeploy_machine.foomac"
	const tfMachineName = "octo-terra-test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMachineBasic(tfMachineName),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployMachineExists(tfVarPrefix),
					resource.TestCheckResourceAttr(
						tfVarPrefix, "name", tfMachineName),
				),
			},
		},
	})
}

func testMachineBasic(machineName string) string {
	config := fmt.Sprintf(`
	data "octopusdeploy_machinepolicy" "default" {
		name = "Default Machine Policy"
	}

	resource "octopusdeploy_environment" "tf_test_env" {
		name           = "OctopusTestMachineBasic"
		description    = "Environment for testing Octopus Machines"
		use_guided_failure = "false"
	}

	resource "octopusdeploy_machine" "foomac" {
		name                            = "%s"
		environments                    = ["${octopusdeploy_environment.tf_test_env.id}"]
		isdisabled                      = true
		machinepolicy                   = "${data.octopusdeploy_machinepolicy.default.id}"
		roles                           = ["Prod"]
		tenanteddeploymentparticipation = "Untenanted"

		endpoint {
		  communicationstyle = "None"
		  thumbprint         = ""
		  uri                = ""
		}
	  }
		`, machineName,
	)
	fmt.Println(config)
	return config
}

func testOctopusDeployMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		return existsMachineHelper(s, client)
	}
}

func existsMachineHelper(s *terraform.State, client *client.Client) error {
	macID := s.RootModule().Resources["octopusdeploy_machine.foomac"].Primary.ID

	if _, err := client.Machines.Get(macID); err != nil {
		return fmt.Errorf("Received an error retrieving machine %s", err)
	}

	return nil
}

func testOctopusDeployMachineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	return destroyMachineHelper(s, client)
}

func destroyMachineHelper(s *terraform.State, apiClient *client.Client) error {
	macID := s.RootModule().Resources["octopusdeploy_machine.foomac"].Primary.ID

	if err := apiClient.Machines.Delete(macID); err != nil {
		if err == client.ErrItemNotFound {
			return nil
		}
		return fmt.Errorf("Received an error retrieving machine %s", err)
	}
	return fmt.Errorf("Machine still exists")
}
