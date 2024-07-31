package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestExpandRunScriptAction(t *testing.T) {
	runScriptAction := expandRunScriptAction(nil)
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"name": nil,
	})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"action_type": nil,
	})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"name": acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	})
	require.NotNil(t, runScriptAction)
}

func (suite *IntegrationTestSuite) TestAccRunScriptAction() {
	feedLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedURI := "http://test.com"
	feedIsEnhancedMode := true
	feedUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccDeploymentProcessCheckDestroy,
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccEnvironmentCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRunScriptAction(),
				),
				Config: testAccRunScriptAction(feedLocalName, feedName, feedURI, feedUsername, feedPassword, feedIsEnhancedMode),
			},
		},
	})
}

func testAccRunScriptAction(feedLocalName string, feedName string, feedURI string, feedUsername string, feedPassword string, feedIsEnhancedMode bool) string {
	return fmt.Sprintf(testAccNuGetFeed(feedLocalName, feedName, feedURI, feedUsername, feedPassword, feedIsEnhancedMode)+"\n"+
		testAccBuildTestAction(`
		run_script_action {
			name                           = "Run Script"
			sort_order = 1
			run_on_server                  = true
			script_file_name               = "Test.ps1"
			script_parameters              = "-Test 1"
			script_source                  = "Package"
			variable_substitution_in_files = "test.json"

			package {
				acquisition_location      = "Server"
				extract_during_deployment = false
				feed_id                   = "${octopusdeploy_nuget_feed.%s.id}"
				name                      = "package2"
				package_id                = "package2"
			}

			primary_package {
				feed_id    = "feeds-builtin"
				package_id = "MyPackage"
			}
		}
	`), feedLocalName)
}

func testAccCheckRunScriptAction() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		process, err := getDeploymentProcess(s, octoClient)
		if err != nil {
			return err
		}

		action := process.Steps[0].Actions[0]

		if action.ActionType != "Octopus.Script" {
			return fmt.Errorf("Action type is incorrect: %s", action.ActionType)
		}

		if action.Properties["Octopus.Action.Script.ScriptFileName"].Value != "Test.ps1" {
			return fmt.Errorf("ScriptFileName is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptFileName"].Value)
		}

		if action.Properties["Octopus.Action.Script.ScriptParameters"].Value != "-Test 1" {
			return fmt.Errorf("ScriptSource is incorrect: %s", action.Properties["Octopus.Action.Script.ScriptParameters"].Value)
		}

		if action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"].Value != "test.json" {
			return fmt.Errorf("TargetFiles is incorrect: %s", action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"].Value)
		}

		return nil
	}
}

func testAccNuGetFeed(localName string, name string, feedURI string, username string, password string, isEnhancedMode bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_nuget_feed" "%s" {
		feed_uri         = "%s"
		is_enhanced_mode = %v
		name             = "%s"
		password         = "%s"
		username         = "%s"
	}`, localName, feedURI, isEnhancedMode, name, password, username)
}
