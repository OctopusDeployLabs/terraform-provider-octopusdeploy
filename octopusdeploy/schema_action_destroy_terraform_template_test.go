package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDestroyTerraformAction(t *testing.T) {
	allowPluginDownloads := acctest.RandIntRange(0, 2) == 0
	applyParameters := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	initParameters := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	pluginCacheDirectory := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	runOnServer := acctest.RandIntRange(0, 2) == 0
	isPlan := acctest.RandIntRange(0, 2) == 0
	workspace := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	scriptSource := "Inline"
	if acctest.RandIntRange(0, 2) == 0 {
		scriptSource = "Package"
	}

	parameters := ""
	source := ""
	if scriptSource == "Inline" {
		variableName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
		variableValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
		source = fmt.Sprintf(`variable \"%s\" { type = string }`, variableName)
		parameters = fmt.Sprintf(`{\"%s\":\"%s\"}`, variableName, variableValue)
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDestroyTerraformAction(name, runOnServer, isPlan, scriptSource, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace, source, parameters),
				),
				Config: testAccDestroyTerraformAction(name, runOnServer, isPlan, scriptSource, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace, source, parameters),
			},
		},
	})
}

func testAccDestroyTerraformAction(name string, runOnServer bool, isPlan bool, templateSource string, allowPluginDownloads bool, applyParameters string, initParameters string, pluginCacheDirectory string, workspace string, template string, templateParameters string) string {
	return testAccBuildTestAction(fmt.Sprintf(`
		destroy_terraform_template_action {
			name          = "%s"
			run_on_server = %v
			is_plan       = %v

			template {
				additional_variable_files         = "additional-variable-files"
				directory                         = "template-directory"
				run_automatic_file_substitution   = false
				target_files                      = "target-files"
			}

			advanced_options {
				allow_additional_plugin_downloads = %v
				apply_parameters                  = "%s"
				init_parameters                   = "%s"
				plugin_cache_directory            = "%s"
				workspace                         = "%s"
			}

			aws_account {
				region            = "us-east-1"
				variable          = "foo"
				use_instance_role = true

				role {
					arn               = "arn"
					external_id       = "external-id"
					role_session_name = "role-session-name"
					session_duration  = 1800
				}
			}

			azure_account {
				variable = "qwe"
			}

			primary_package {
				package_id = "MyPackage"
				feed_id = "feeds-builtin"
			}
    }`, name, runOnServer, isPlan, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace))
}

func testAccCheckDestroyTerraformAction(name string, runOnServer bool, isPlan bool, scriptSource string, allowPluginDownloads bool, applyParameters string, initParameters string, pluginCacheDirectory string, workspace string, source string, parameters string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if !isPlan && action.ActionType != "Octopus.TerraformDestroy" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if isPlan && action.ActionType != "Octopus.TerraformPlanDestroy" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Terraform.AdditionalInitParams"].Value != initParameters {
			return fmt.Errorf("AdditionalInitParams: %s", action.Properties["Octopus.Action.Terraform.AdditionalInitParams"].Value)
		}

		if v, _ := strconv.ParseBool(action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"].Value); v != allowPluginDownloads {
			return fmt.Errorf("AllowPluginDownloads: %s", action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"].Value)
		}

		return nil
	}
}
