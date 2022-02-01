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

func TestAccOctopusDeployApplyTerraformAction(t *testing.T) {
	allowPluginDownloads := acctest.RandIntRange(0, 2) == 0
	applyParameters := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	initParameters := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	pluginCacheDirectory := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	runOnServer := acctest.RandIntRange(0, 2) == 0
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
					testAccCheckApplyTerraformAction(name, runOnServer, scriptSource, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace, source, parameters),
				),
				Config: testAccApplyTerraformAction(name, runOnServer, scriptSource, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace, source, parameters),
			},
		},
	})
}

func testAccApplyTerraformAction(name string, runOnServer bool, templateSource string, allowPluginDownloads bool, applyParameters string, initParameters string, pluginCacheDirectory string, workspace string, template string, templateParameters string) string {
	return testAccBuildTestAction(fmt.Sprintf(`
		apply_terraform_template_action {
			name          = "%s"
			run_on_server = %v

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

			google_cloud_account {
				variable = "boo"
				use_vm_service_account = true
				project = "projectboo"
				region = "regionboo"
				zone = "zoneboo"
				service_account_emails = "my@octopus.com"
				impersonate_service_account = true
			}

			primary_package {
				package_id = "MyPackage"
				feed_id = "feeds-builtin"
			}
    }`, name, runOnServer, allowPluginDownloads, applyParameters, initParameters, pluginCacheDirectory, workspace))
}

func testAccCheckApplyTerraformAction(name string, runOnServer bool, scriptSource string, allowPluginDownloads bool, applyParameters string, initParameters string, pluginCacheDirectory string, workspace string, source string, parameters string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)

		process, err := getDeploymentProcess(s, client)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.TerraformApply" {
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
