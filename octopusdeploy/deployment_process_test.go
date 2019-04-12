package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployDeploymentProcessBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentProcessBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployDeploymentProcess(),
				),
			},
		},
	})
}

func testAccDeploymentProcessBasic() string {
	return `
		resource "octopusdeploy_lifecycle" "test" {
			name = "Test Lifecycle"
		}

		resource "octopusdeploy_project_group" "test" {
			name = "Test Group"
		}

		resource "octopusdeploy_project" "test" {
			name             = "Test Project"
			lifecycle_id     = "${octopusdeploy_lifecycle.test.id}"
			project_group_id = "${octopusdeploy_project_group.test.id}"
		}

		resource "octopusdeploy_deployment_process" "test" {
			project_id = "${octopusdeploy_project.test.id}"

			step {
				name = "Test"
				target_roles = ["A", "B"]
				package_requirement = "AfterPackageAcquisition"
				condition = "Variable"
				condition_expression = "#{run}"
				start_trigger = "StartWithPrevious"
				window_size = "5"
				

				action {
					name = "Test"
					action_type = "Octopus.Script"
					disabled = false
					required = true
					worker_pool_id = "WorkerPools-1"
					environments = ["Environments-1"]
					//excluded_environments = ["Environments-2"]
					//channels = ["Channels-1"]
					//tenant_tags = ["tag/tag"]
					
					primary_package {
						package_id = "MyPackage"
						feed_id = "feeds-builtin"
						acquisition_location = "ExecutionTarget"
					}

					package {
						name = "ThePackage"
						package_id = "MyPackage"
						feed_id = "feeds-builtin"
						acquisition_location = "NotAcquired"
						extract_during_deployment = true

						property {
							key = "WhatIsThis"
							value = "Dunno"
						}

					}

					package {
						name = "ThePackage2"
						package_id = "MyPackage2"
						feed_id = "feeds-builtin"
						acquisition_location = "NotAcquired"
						extract_during_deployment = true

						property {
							key = "WhatIsThis"
							value = "Dunno"
						}
					}

					property {
						key = "Octopus.Action.Script.ScriptFileName"
						value = "Run.ps132"
					}

					property {
						key = "Octopus.Action.Script.ScriptSource"
						value = "Package"
					}

				}
			}

 			step {
 			       name = "Step2"
 			       start_trigger = "StartWithPrevious"
			
 			       action {
 			           name = "Step2"
 			           action_type = "Octopus.Script"
 			           run_on_server = true
			
 			           property {
 			               key = "Octopus.Action.Script.ScriptBody"
 			               value = "Write-Host 'hi'"
 			           }
 			       }
			} 
		}
		`
}

func testAccBuildTestAction(action string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "test" {
			name = "Test Lifecycle"
		}

		resource "octopusdeploy_project_group" "test" {
			name = "Test Group"
		}

		resource "octopusdeploy_project" "test" {
			name             = "Test Project"
			lifecycle_id     = "${octopusdeploy_lifecycle.test.id}"
			project_group_id = "${octopusdeploy_project_group.test.id}"
		}

		resource "octopusdeploy_deployment_process" "test" {
			project_id = "${octopusdeploy_project.test.id}"

			step {
				name = "Test"
				target_roles = ["WebServer"]

				%s
			}
		}
		`, action)
}

func testAccCheckOctopusDeployDeploymentProcessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyProjectHelper(s, client); err != nil {
		return err
	}
	if err := destroyHelperProjectGroup(s, client); err != nil {
		return err
	}
	if err := destroyHelperLifecycle(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployDeploymentProcess() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		expectedNumberOfSteps := 2
		numberOfSteps := len(process.Steps)
		if numberOfSteps != expectedNumberOfSteps {
			return fmt.Errorf("Deployment process has %d steps instead of the expected %d", numberOfSteps, expectedNumberOfSteps)
		}

		if process.Steps[0].Actions[0].Properties["Octopus.Action.RunOnServer"] != "true" {
			return fmt.Errorf("The RunOnServer property has not been set to true on the deployment process")
		}

		return nil
	}
}

func getDeploymentProcess(s *terraform.State, client *octopusdeploy.Client) (*octopusdeploy.DeploymentProcess, error) {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_deployment_process" {
			return client.DeploymentProcess.Get(r.Primary.ID)
		}
	}
	return nil, fmt.Errorf("No deployment process found in the terraform resources")
}
