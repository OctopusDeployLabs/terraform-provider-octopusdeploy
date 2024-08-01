package octopusdeploy

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccOctopusDeployDeployWindowsServiceAction() {
	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeployWindowsServiceAction(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployWindowsServiceActionOrFeature(suite, "Octopus.WindowsService"),
				),
			},
		},
	})
}

func (suite *IntegrationTestSuite) TestAccOctopusDeployWindowsServiceFeature() {
	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsServiceFeature(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployWindowsServiceActionOrFeature(suite, "Octopus.TentaclePackage"),
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
			sort_order = 1

			primary_package {
				package_id = "MyPackage"
			}
		}
	`)
}

func testAccWindowsServiceFeature() string {
	return testAccBuildTestAction(`
		deploy_package_action {
			name = "Test"
			sort_order = 1

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

func testAccCheckDeployWindowsServiceActionOrFeature(suite *IntegrationTestSuite, expectedActionType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		process, err := getDeploymentProcess(s, suite.octoClient)
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

		if action.Properties["Octopus.Action.WindowsService.ServiceName"].Value != "MyService" {
			return fmt.Errorf("Service Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceName"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.DisplayName"].Value != "My Service" {
			return fmt.Errorf("Display Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.DisplayName"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.Description"].Value != "Do stuff" {
			return fmt.Errorf("Description is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Description"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.ExecutablePath"].Value != "MyService.exe" {
			return fmt.Errorf("Executable Path is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ExecutablePath"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.Arguments"].Value != "-arg" {
			return fmt.Errorf("Arguments is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Arguments"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.ServiceAccount"].Value != "_CUSTOM" {
			return fmt.Errorf("Service Account is incorrect: %s", action.Properties["Octopus.Action.WindowsService.ServiceAccount"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountName"].Value != "User" {
			return fmt.Errorf("Custom Account Name is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountName"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"].Value != "Password" {
			return fmt.Errorf("Custom Account Password is incorrect: %s", action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.StartMode"].Value != "manual" {
			return fmt.Errorf("Start Mode is incorrect: %s", action.Properties["Octopus.Action.WindowsService.StartMode"].Value)
		}

		if action.Properties["Octopus.Action.WindowsService.Dependencies"].Value != "OtherService" {
			return fmt.Errorf("Dependencies is incorrect: %s", action.Properties["Octopus.Action.WindowsService.Dependencies"].Value)
		}

		return nil
	}
}
