package octopusdeploy

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployProjectBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const allowDeploymentsToNoTargets = "true"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "allow_deployments_to_no_targets", allowDeploymentsToNoTargets),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectWithDeploymentStepWindowsService(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const serviceName = "Epic Service"
	const executablePath = `bin\\MyService.exe` // needs 4 slashes to appear in the TF config as a double slash
	const stepName = "Deploying Epic Service"
	const packageName = "MyPackage"
	targetRoles := []string{"Lab1", "Lab2"}
	projectIDRegex := regexp.MustCompile(`Projects\-`)
	deploymentProcessIDRegex := regexp.MustCompile(`deploymentprocess\-Projects\-.*`)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWithDeploymentStepWindowsService(projectName, lifeCycleID, serviceName, executablePath, stepName, packageName, targetRoles),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestMatchResourceAttr(
						terraformNamePrefix, "id", projectIDRegex),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.service_name", serviceName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.step_name", stepName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.target_roles.0", targetRoles[0]),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.target_roles.1", targetRoles[1]),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.executable_path", strings.Replace(executablePath, "\\\\", "\\", 1)), // need to scape the backslashes
					resource.TestMatchResourceAttr(
						terraformNamePrefix, "deployment_process_id", deploymentProcessIDRegex),
				),
			},
		},
	})
}

//nolint:govet
func TestAccOctopusDeployProjectWithUpdate(t *testing.T) {
	return

	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const allowDeploymentsToNoTargets = "true"
	inlineScriptRegex := regexp.MustCompile(`.*Get\-Process.*`)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			// create project with no description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
				),
			},
			// create update it with a description + build steps
			{
				Config: testAccWithMultipleDeploymentStepWindowsService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", "Project Name"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", "Lifecycles-1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", "My Awesome Description"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.service_name", "My First Service"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.step_name", "Deploy My First Service"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.target_roles.0", "Role1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.target_roles.1", "Role2"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.executable_path", "C:\\MyService\\my_service.exe"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.service_name", "My Second Service"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.step_name", "Deploy My Second Service"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.target_roles.0", "Role3"),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.target_roles.1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.executable_path", "C:\\MyService\\my_service2.exe"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.configuration_transforms", "false"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.1.windows_service.1.configuration_variables", "false"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.step_name", "Run Cleanup Script"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.target_roles.0", "MyRole1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.target_roles.1", "MyRole2"),
					resource.TestMatchResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.script_body", inlineScriptRegex),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.script_type", "PowerShell"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.2.inline_script.0.step_condition", "success"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.step_name", "Deploy Website"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.target_roles.0", "MyRole1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.website_name", "Awesome Website"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.application_pool_name", "MyAppPool"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.application_pool_framework", "v2.0"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.step_condition", "failure"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.basic_authentication", "true"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.anonymous_authentication", "true"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.3.iis_website.0.json_file_variable_replacement", "appsettings.json,Config\\*.json"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.feed_id", "feeds-builtin"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.package", "cleanup.yolo"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.script_file_name", "bin\\cleanup.ps1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.script_parameters", "-Force"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.step_name", "Run Verify From Package Script"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.step_condition", "success"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.target_roles.0", "MyRole1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "deployment_step.4.package_script.0.target_roles.1", "MyRole2"),
				),
			},
			// update again by remove its description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "lifecycle_id", lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", ""),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.iis_website.0.step_name"),
				),
			},
		},
	})
}

func testAccProjectBasic(name, lifeCycleID, allowDeploymentsToNoTargets string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			name           = "%s"
			lifecycle_id    = "%s"
			project_group_id = "${octopusdeploy_project_group.foo.id}"
			allow_deployments_to_no_targets = "%s"
		}
		`,
		name, lifeCycleID, allowDeploymentsToNoTargets,
	)
}

const testAccWithMultipleDeploymentStepWindowsService = `
resource "octopusdeploy_project_group" "foo" {
	name = "Integration Test Project Group"
}

resource "octopusdeploy_project" "foo" {
	name             = "Project Name"
	lifecycle_id     = "Lifecycles-1"
	project_group_id = "${octopusdeploy_project_group.foo.id}"
	description      = "My Awesome Description"

	deployment_step {
		windows_service {
			executable_path          = "C:\\MyService\\my_service.exe"
			package                  = "MyPackage"
			service_name             = "My First Service"
			step_name                = "Deploy My First Service"

			target_roles = [
			"Role1",
			"Role2"
			]
		}
	}

	deployment_step {
		windows_service {
			configuration_transforms = false
			configuration_variables  = false
			executable_path          = "C:\\MyService\\my_service2.exe"
			package                  = "MyServicePackage"
			service_account          = "NewServiceAccount"
			service_name             = "My Second Service"
			service_start_mode       = "demand"
			step_name                = "Deploy My Second Service"
			step_start_trigger       = "StartWithPrevious"

			target_roles = [
				"Role3",
			]
		}
	}

	deployment_step {
		inline_script {
			step_name   = "Run Cleanup Script"
			script_type = "PowerShell"

			script_body = <<EOF
		$x = Get-Process
		foreach ($p in $x) {
			Write-Output $p.Name
		}
		EOF

			target_roles = [
				"MyRole1",
				"MyRole2",
			]
		}
	}

	deployment_step {
		iis_website {
			anonymous_authentication       = true
			application_pool_framework     = "v2.0"
			application_pool_name          = "MyAppPool"
			basic_authentication           = true
			json_file_variable_replacement = "appsettings.json,Config\\*.json"
			package                        = "MyWebsitePackage"
			step_condition                 = "failure"
			step_name                      = "Deploy Website"
			website_name                   = "Awesome Website"

			target_roles = [
				"MyRole1",
			]
		}
	}

	deployment_step {
		package_script {
			step_name         = "Run Verify From Package Script"
			package           = "cleanup.yolo"
			script_file_name  = "bin\\cleanup.ps1"
			script_parameters = "-Force"

			target_roles = [
			"MyRole1",
			"MyRole2",
			]
		}
	}
}
`

func testAccWithDeploymentStepWindowsService(name, lifeCycleID, serviceName, executablePath, stepName, packageName string, targetRoles []string) string {
	return fmt.Sprintf(`

		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			name             = "%s"
			lifecycle_id     = "%s"
			project_group_id = "${octopusdeploy_project_group.foo.id}"

			deployment_step {
				windows_service {
					executable_path = "%s"
					service_name    = "%s"
					step_name       = "%s"
					package         = "%s"

					target_roles = [
						"%s",
					]
				}
			}
		}
		`,
		name, lifeCycleID, executablePath, serviceName, stepName, packageName, strings.Join(targetRoles, "\",\""),
	)
}

func testAccCheckOctopusDeployProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyProjectHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyProjectHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Project.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving project %s", err)
		}
		return fmt.Errorf("Project still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopus_deploy_project" {
			if _, err := client.Project.Get(r.Primary.ID); err != nil {
				return fmt.Errorf("Received an error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
