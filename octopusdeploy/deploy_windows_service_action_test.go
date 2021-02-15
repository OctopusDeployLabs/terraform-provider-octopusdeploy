package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDeployWindowsServiceAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployWindowsServiceAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployWindowsServiceActionOrFeature("Octopus.WindowsService"),
				),
			},
		},
	})
}

func TestAccOctopusDeployWindowsServiceFeature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsServiceFeature(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployWindowsServiceActionOrFeature("Octopus.TentaclePackage"),
				),
			},
		},
	})
}

func testAccDeployWindowsServiceAction() string {
	return testAccBuildTestAction(`
		deploy_windows_service_action {
			arguments = "-arg"
			custom_account_name = "User"
			custom_account_password = "Password"
			dependencies = "OtherService"
			description = "Do stuff"
			display_name = "My Service"
			executable_path = "MyService.exe"
			name = "Test"
			service_account = "_CUSTOM"
			service_name = "MyService"
			start_mode = "manual"

			package {
				package_id = "MyPackage"
			}
		}
	`)
}

func testAccWindowsServiceFeature() string {
	return testAccBuildTestAction(`
		deploy_package_action {
			name = "Test"

			primary_package {
				package_id = "MyPackage"
			}

			windows_service {
				arguments = "-arg"
				custom_account_name = "User"
				custom_account_password = "Password"
				description = "Do stuff"
				dependencies = "OtherService"
				display_name = "My Service"
				executable_path = "MyService.exe"
				service_account = "_CUSTOM"
				service_name = "MyService"
				start_mode = "manual"
			}
		}
	`)
}

func testAccCheckDeployWindowsServiceActionOrFeature(expectedActionType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != expectedActionType {
			return fmt.Errorf("Action type is incorrect: %s, expected: %s", action.ActionType, expectedActionType)
		}

		if len(action.Packages) == 0 {
			return fmt.Errorf("No package")
		}

		// if action.Properties["Octopus.Action.WindowsService.CreateOrUpdateService"] != "True" {
		// 	return fmt.Errorf("Windows Service feature is not enabled")
		// }

		if action.Properties["Octopus.Action.WindowsService.ServiceName"] != "MyService" {
			return fmt.Errorf("Service Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceName"])
		}

		if action.Properties["Octopus.Action.WindowsService.DisplayName"] != "My Service" {
			return fmt.Errorf("Display Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.DisplayName"])
		}

		if action.Properties["Octopus.Action.WindowsService.Description"] != "Do stuff" {
			return fmt.Errorf("Description is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Description"])
		}

		if action.Properties["Octopus.Action.WindowsService.ExecutablePath"] != "MyService.exe" {
			return fmt.Errorf("Executable Path is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ExecutablePath"])
		}

		if action.Properties["Octopus.Action.WindowsService.Arguments"] != "-arg" {
			return fmt.Errorf("Arguments is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Arguments"])
		}

		if action.Properties["Octopus.Action.WindowsService.ServiceAccount"] != "_CUSTOM" {
			return fmt.Errorf("Service Account is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceAccount"])
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountName"] != "User" {
			return fmt.Errorf("Custom Account Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountName"])
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"] != "Password" {
			return fmt.Errorf("Custom Account Password is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"])
		}

		if action.Properties["Octopus.Action.WindowsService.StartMode"] != "manual" {
			return fmt.Errorf("Start Mode is incorrect: %s", action.Properties["Octopus.Action.WindowsService.StartMode"])
		}

		if action.Properties["Octopus.Action.WindowsService.Dependencies"] != "OtherService" {
			return fmt.Errorf("Dependencies is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Dependencies"])
		}

		return nil
	}
}
