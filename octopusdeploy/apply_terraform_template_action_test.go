package octopusdeploy

// import (
// 	"fmt"
// 	"strconv"
// 	"testing"

// 	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// func TestAccOctopusDeployApplyTerraformAction(t *testing.T) {
// 	additionalInitialParameters := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	allowPluginDownloads := acctest.RandIntRange(0, 2) == 0

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckOctopusDeployDeploymentProcessDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckApplyTerraformAction(additionalInitialParameters, allowPluginDownloads),
// 				),
// 				Config: testAccApplyTerraformAction(additionalInitialParameters, allowPluginDownloads),
// 			},
// 		},
// 	})
// }

// func testAccApplyTerraformAction(additionalInitParams string, allowPluginDownloads bool) string {
// 	return testAccBuildTestAction(fmt.Sprintf(`
// 		apply_terraform_template_action {
// 			name = "%s"
// 			run_on_server = %v

// 			template {
// 				template_source = "%s"
// 				template_parameters = "%s"
// 			}

// 			managed_accounts {
// 				enable_aws_account_integration = %v
// 				enable_azure_account_integration = %v
// 				managed_account = "%s"
// 			}

// 			advanced_options {
// 				allow_additional_plugin_downloads = %v
// 				apply_parameters = "%s"
// 				init_parameters = "%s"
// 				plugin_cache_directory = "%s"
// 				run_automatic_file_substitution = %v
// 				workspace = "%s"
// 			}

// 			primary_package {
// 				package_id = "MyPackage"
// 				feed_id = "feeds-builtin"
// 			}
//     }`, additionalInitParams, allowPluginDownloads, enableAwsAccountIntegration, enableAzureAccountIntegration, managedAccount, name, pluginsDirectory, runAutomaticFileSubstitution, runOnServer, terraformTemplate, terraformTemplateParameters, terraformWorkspace))
// }

// func testAccCheckApplyTerraformAction(additionalInitialParameters string, allowPluginDownloads bool) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		client := testAccProvider.Meta().(*octopusdeploy.Client)

// 		process, err := getDeploymentProcess(s, client)
// 		if err != nil {
// 			return err
// 		}

// 		action := process.Steps[0].Actions[0]

// 		if action.ActionType != "Octopus.TerraformApply" {
// 			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
// 		}

// 		if action.Properties["Octopus.Action.Terraform.AdditionalInitParams"].Value != additionalInitialParameters {
// 			return fmt.Errorf("AdditionalInitParams: %s", action.Properties["Octopus.Action.Terraform.AdditionalInitParams"].Value)
// 		}

// 		if v, _ := strconv.ParseBool(action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"].Value); v != allowPluginDownloads {
// 			return fmt.Errorf("AllowPluginDownloads: %s", action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"].Value)
// 		}

// 		return nil
// 	}
// }
