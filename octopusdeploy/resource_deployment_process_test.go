package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDeploymentProcessWithPackage(t *testing.T) {
	actionName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	packageID := "Octopus.Cli"
	packageName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	stepName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_deployment_process." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccDeploymentProcessCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDeploymentProcessExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
				Config: testAccDeploymentProcessWithPackage(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, stepName, actionName, packageName, packageID),
			},
		},
	})
}

func testAccDeploymentProcessWithPackage(localName string, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, stepName string, actionName string, packageName string, packageID string) string {
	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription)+"\n"+`
		resource "octopusdeploy_deployment_process" "%s" {
			project_id = "${octopusdeploy_project.%s.id}"

			step {
				name         = "%s"
				target_roles = ["Foo"]

				deploy_tentacle_package_action {
					name         = "%s"

					package {
						name       = "%s"
						package_id = "%s"
					}
				}
			}
		}`, localName, projectLocalName, stepName, actionName, packageName, packageID)
}

func TestAccOctopusDeployDeploymentProcessBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_process." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccDeploymentProcessCheckDestroy,
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDeploymentProcessExists(resourceName),
					testAccCheckOctopusDeployDeploymentProcess(),
				),
				Config: testAccDeploymentProcessBasic(localName),
			},
		},
	})
}

func TestAccOctopusDeployDeploymentProcessWithActionTemplate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_process." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccDeploymentProcessCheckDestroy,
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccDeploymentProcessExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "step.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "step.0.condition", "Success"),
					resource.TestCheckResourceAttr(resourceName, "step.0.name", "Terraform - PlanV2"),
					resource.TestCheckResourceAttr(resourceName, "step.0.start_trigger", "StartAfterPrevious"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.action_type", "Octopus.TerraformPlan"),
					// resource.TestCheckResourceAttr(prefix, "step.0.action.0.can_be_used_for_project_versioning", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.channels"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.container"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.condition", "Success"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.environments"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.run_on_server", "true"),
				),
				Config: testAccProcessWithActionTemplate(localName),
			},
		},
	})
}

func testAccDeploymentProcessBasic(localName string) string {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription)+"\n"+
		`resource "octopusdeploy_deployment_process" "%s" {
			branch     = "%s"
			project_id = "${octopusdeploy_project.%s.id}"

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
		}`, localName, projectLocalName)
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

				// template {
				//   community_action_template_id = "CommunityActionTemplates-27"
				//   id                           = "ActionTemplates-281"
				//   version                      = 3
				// }
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
				name         = "Test"
				target_roles = ["Foo", "Bar", "Quux"]

				%s
			}
		}`, projectID, action)
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

func testAccDeploymentProcessExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if _, err := client.DeploymentProcesses.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccDeploymentProcessCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_deployment_process" {
			continue
		}

		if deploymentProcess, err := client.DeploymentProcesses.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("deployment process (%s) still exists", deploymentProcess.GetID())
		}
	}

	return nil
}
