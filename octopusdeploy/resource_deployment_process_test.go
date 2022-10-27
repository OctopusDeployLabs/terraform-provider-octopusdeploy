package octopusdeploy

import (
	"fmt"
	"os"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// func TestAccDeploymentProcess(t *testing.T) {
// 	options := test.NewDeploymentProcessTestOptions()
// 	resourceName := "octopusdeploy_deployment_process." + options.LocalName

// 	resource.Test(t,resource.TestCase{
// 		CheckDestroy: testAccDeploymentProcessCheckDestroy,
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccDeploymentProcessExists(resourceName),
// 					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
// 					resource.TestCheckResourceAttr(resourceName, "step.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "step.0.name", options.StepName),
// 					resource.TestCheckResourceAttr(resourceName, "step.0.action.#", "1"),
// 					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.action_type", options.ActionType),
// 					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.name", options.ActionName),
// 				),
// 				Config: testAccDeploymentProcessWithOptions(options),
// 			},
// 		},
// 	})
// }

// func testAccDeploymentProcessWithOptions(options *test.DeploymentProcessTestOptions) string {
// 	return fmt.Sprintf(testAccProjectWithOptions(options.ProjectCreateTestOptions())+"\n"+`
// 		resource "octopusdeploy_deployment_process" "%s" {
// 			project_id = octopusdeploy_project.%s.id

// 			step {
// 				name         = "%s"
// 				target_roles = ["Foo"]

// 				action {
// 					action_type = "%s"
// 					name        = "%s"

// 					package {
// 						name       = "%s"
// 						package_id = "%s"
// 					}
// 				}
// 			}
// 		}`, options.LocalName, options.Project.LocalName, options.StepName, options.ActionType, options.ActionName, options.PackageName, options.PackageID)
// }

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
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "step.#", "2"),
				),
				Config: testAccDeploymentProcessBasic(localName),
			},
		},
	})
}

func TestAccOctopusDeployDeploymentProcessWithActionTemplate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_process." + localName

	// TODO: replace with client reference
	spaceID := os.Getenv("OCTOPUS_SPACE")

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
				Config: testAccProcessWithActionTemplate(spaceID, localName),
			},
		},
	})
}

func TestAccOctopusDeployDeploymentProcessWithImpliedPrimaryPackage(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_process." + localName

	// TODO: replace with client reference
	spaceID := os.Getenv("OCTOPUS_SPACE")

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
					resource.TestCheckResourceAttr(resourceName, "step.0.name", "Transfer DG Package"),
					resource.TestCheckResourceAttr(resourceName, "step.0.start_trigger", "StartAfterPrevious"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.action_type", "Octopus.TransferPackage"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.channels"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.container"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.condition", "Success"),
					resource.TestCheckNoResourceAttr(resourceName, "step.0.action.0.environments"),
					resource.TestCheckResourceAttr(resourceName, "step.0.action.0.run_on_server", "false"),
				),
				Config: testAccProcessWithImpliedPrimaryPackage(spaceID, localName),
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
			project_id = octopusdeploy_project.%s.id

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

func testAccProcessWithImpliedPrimaryPackage(spaceID string, localName string) string {
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
				condition           = "Success"
				name                = "Transfer DG Package"
				package_requirement = "LetOctopusDecide"
				start_trigger       = "StartAfterPrevious"
				target_roles        = ["test-1"]
				properties          = {
					"Octopus.Action.TargetRoles": "test-1"
				}

				action {
					action_type                        = "Octopus.TransferPackage"
					can_be_used_for_project_versioning = true
					condition                          = "Success"
					is_disabled                        = false
					is_required                        = false
					name                               = "Test"

					primary_package {
						acquisition_location = "Server"
						feed_id              = "feeds-builtin"
						package_id           = "test-package"

						properties = {
							"SelectionMode" = "immediate"
						}
					}

					properties = {
							"Octopus.Action.Package.TransferPath" = "/var"
							"Octopus.Action.Package.PackageId" = "test-package"
							"Octopus.Action.Package.DownloadOnTentacle" = "False"
					}
				}
			}
  		}`, localName, projectID)
}

func testAccProcessWithActionTemplate(spaceID string, localName string) string {
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
			  condition     = "Success"
			  name          = "Terraform - PlanV2"
			  start_trigger = "StartAfterPrevious"

			  run_script_action {
			    can_be_used_for_project_versioning = true
				is_disabled                        = false
				is_required                        = true
				name                               = "AWS - Create a Security Group"
				run_on_server                      = true
				script_body                        = "// hello"

				template {
				  community_action_template_id = "CommunityActionTemplates-27"
				  id                           = "ActionTemplates-281"
				  version                      = 3
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
				name         = "Test"
				target_roles = ["Foo", "Bar", "Quux"]

				%s
			}
		}`, projectID, action)
}

func testAccCheckOctopusDeployDeploymentProcess() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)

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

func getDeploymentProcess(s *terraform.State, client *client.Client) (*deployments.DeploymentProcess, error) {
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

		client := testAccProvider.Meta().(*client.Client)
		if _, err := client.DeploymentProcesses.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccDeploymentProcessCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
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
