package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDeploymentProcessBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

func TestAccOctopusDeployDeploymentProcessWithActionTemplate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_deployment_process." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttr(prefix, "step.#", "1"),
					resource.TestCheckResourceAttr(prefix, "step.0.condition", "Success"),
					resource.TestCheckResourceAttr(prefix, "step.0.name", "Terraform - PlanV2"),
					resource.TestCheckResourceAttr(prefix, "step.0.start_trigger", "StartAfterPrevious"),
					resource.TestCheckResourceAttr(prefix, "step.0.action.#", "1"),
					resource.TestCheckResourceAttr(prefix, "step.0.action.0.action_type", "Octopus.TerraformPlan"),
					// resource.TestCheckResourceAttr(prefix, "step.0.action.0.can_be_used_for_project_versioning", "true"),
					resource.TestCheckNoResourceAttr(prefix, "step.0.action.0.channels"),
					resource.TestCheckNoResourceAttr(prefix, "step.0.action.0.container"),
					resource.TestCheckResourceAttr(prefix, "step.0.action.0.condition", "Success"),
					resource.TestCheckNoResourceAttr(prefix, "step.0.action.0.environments"),
					resource.TestCheckResourceAttr(prefix, "step.0.action.0.run_on_server", "true"),
				),
				Config: testAccProcessWithActionTemplate(localName),
			},
		},
	})
}

func testAccDeploymentProcessBasic() string {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	projectID := "octopusdeploy_project." + localName + ".id"

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description)+"\n"+
		`resource "octopusdeploy_deployment_process" "test" {
			project_id = %s

			step {
				condition = "Variable"
				condition_expression = "#{run}"
				name = "Test"
				package_requirement = "AfterPackageAcquisition"
				start_trigger = "StartAfterPrevious"
				target_roles = ["A", "B"]
				window_size = "5"

				run_script_action {
					// channels = ["Channels-1"]
					// environments = ["Environments-1"]
					// excluded_environments = ["Environments-2"]
					is_disabled = false
					is_required = true
					name = "Test"
					run_on_server = true
					script_file_name = "Run.ps132"
					script_source = "Package"
					tenant_tags = ["tag/tag"]
					
					primary_package {
						acquisition_location = "Server"
						feed_id = "feeds-builtin"
						package_id = "MyPackage"
					}

					package {
						acquisition_location = "NotAcquired"
						extract_during_deployment = true
						feed_id = "feeds-builtin"
						name = "ThePackage"
						package_id = "MyPackage"
					}

					package {
						acquisition_location = "NotAcquired"
						extract_during_deployment = true
						feed_id = "feeds-builtin"
						name = "ThePackage2"
						package_id = "MyPackage2"
					}
				}
			}

 			step {
			  name = "Step2"
			  start_trigger = "StartWithPrevious"
			  target_roles = ["WebServer"]
	
			  run_script_action {
				  name = "Step2"
				  run_on_server = true
				  script_body = "Write-Host 'hi'"
			  }
			} 
		}`, projectID)
}

func testAccProcessWithActionTemplate(localName string) string {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	projectID := "octopusdeploy_project." + projectLocalName + ".id"

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, name, description)+"\n"+
		`resource "octopusdeploy_deployment_process" "%s" {
			project_id = %s

			step {
			  condition = "Success"
			  name = "Terraform - PlanV2"
			  start_trigger = "StartAfterPrevious"

			  action {
			    action_type = "Octopus.Script"
			    can_be_used_for_project_versioning = true
				is_disabled = false
				is_required = true
				name = "AWS - Create a Security Group"
				run_on_server = true

				template {
				  community_action_template_id = "CommunityActionTemplates-27"
				  id = "ActionTemplates-281"
				  version = 3
				}
			  }
			}
		}`, localName, projectID)
}

func testAccBuildTestAction(action string) string {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	projectID := "octopusdeploy_project." + localName + ".id"

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description)+"\n"+
		`resource "octopusdeploy_deployment_process" "test" {
			project_id = %s

			step {
				name = "Test"
				target_roles = ["Foo", "Bar", "Quux"]

				%s
			}
		}`, projectID, action)
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

		if process.Steps[0].Actions[0].Properties["Octopus.Action.RunOnServer"].Value != "true" {
			return fmt.Errorf("The RunOnServer property has not been set to true on the deployment process")
		}

		return nil
	}
}

func getDeploymentProcess(s *terraform.State, client *octopusdeploy.Client) (*octopusdeploy.DeploymentProcess, error) {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_deployment_process" {
			return client.DeploymentProcesses.GetByID(r.Primary.ID)
		}
	}
	return nil, fmt.Errorf("No deployment process found in the terraform resources")
}
