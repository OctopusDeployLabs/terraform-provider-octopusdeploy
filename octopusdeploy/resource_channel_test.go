package octopusdeploy

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployChannelBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_channel." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccChannelCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
				Config: testAccChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, name, description),
			},
		},
	})
}

func TestAccOctopusDeployChannelBasicWithUpdate(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_channel." + localName
	const channelName = "Funky Channel"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccChannelCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// create baseline channel
			{
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channelName),
					resource.TestCheckResourceAttr(resourceName, "description", "this is funky"),
				),
				Config: testAccChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelName, "this is funky"),
			},
			// update channel with a new description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channelName),
					resource.TestCheckResourceAttr(resourceName, "description", "funky it is"),
				),
				Config: testAccChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelName, "funky it is"),
			},
		},
	})
}

func TestAccOctopusDeployChannelWithOneRule(t *testing.T) {
	internalTest.SkipCI(t, "action_package blocks required on rule, this test is out of date.")
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const channelDescription = "this is Funky"
	const actionName = "Funky Action"
	const versionRange = "1.0"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccChannelCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{ // create channel with one rule
				Config: testAccChannelWithOneRule(channelName, channelDescription, versionRange, actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", channelDescription),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.version_range", versionRange),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.actions.0", actionName),
				),
			},
		},
	})
}

func TestAccOctopusDeployChannelWithOneRuleWithUpdate(t *testing.T) {
	internalTest.SkipCI(t, "action_package blocks required on rule, this test is out of date.")
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const updatedChannelName = "Updated Channel"
	const channelDescription = "this is Funky"
	const updatedChannelDescription = "this is updated"
	const versionRange = "1.0"
	const updatedVersionRange = "2.5"
	const actionName = "Funky Action"
	const updatedActionName = "Updated Action"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccChannelCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{ // create baseline channel
				Config: testAccChannelWithOneRule(channelName, channelDescription, versionRange, actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", channelDescription),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.version_range", versionRange),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.actions.0", actionName),
				),
			},
			{ // create updated channel with new values
				Config: testAccChannelWithOneRule(updatedChannelName, updatedChannelDescription, updatedVersionRange, updatedActionName),
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", updatedChannelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", updatedChannelDescription),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.version_range", updatedVersionRange),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.actions.0", updatedActionName),
				),
			},
		},
	})
}

func TestAccOctopusDeployChannelWithTwoRules(t *testing.T) {
	internalTest.SkipCI(t, "action_package blocks required on rule, this test is out of date.")
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const channelDescription = "this is Funky"
	const versionRange1 = "1.0"
	const actionName1 = "Funky Action"
	const versionRange2 = "2.0"
	const actionName2 = "Action-2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			{ // create channel with two rules
				Config: testAccChannelWithTwoRules(channelName, channelDescription, versionRange1, actionName1, versionRange2, actionName2),
				Check: resource.ComposeTestCheckFunc(
					testAccChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", channelDescription),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.version_range", versionRange1),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.actions.0", actionName1),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.1.version_range", versionRange2),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.1.actions.0", actionName2),
				),
			},
		},
	})
}

func testAccChannelBasic(localName string, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, name string, description string) string {
	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription)+"\n"+`
		resource "octopusdeploy_channel" "%s" {
			description = "%s"
			name        = "%s"
			project_id  = octopusdeploy_project.%s.id
		}`, localName, description, name, projectLocalName)
}

func testAccProjectBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	projectGroup := internalTest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internalTest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

			template {
				default_value = "default-value"
				help_text     = "help-test"
				label         = "label"
				name          = "2"

				display_settings = {
					"Octopus.ControlType": "SingleLineText"
				}
			}

			template {
				default_value = "default-value"
				help_text     = "help-test"
				label         = "label"
				name          = "1"

				display_settings = {
					"Octopus.ControlType": "SingleLineText"
				}
			}

		  //   connectivity_policy {
		//     allow_deployments_to_no_targets = true
		// 	skip_machine_behavior           = "None"
		//   }

		//   version_control_settings {
		// 	default_branch = "foo"
		// 	url            = "https://example.com/"
		// 	username       = "bar"
		//   }

		//   versioning_strategy {
		//     template = "alskdjaslkdj"
		//   }
		}`, localName, description, lifecycleLocalName, name, projectGroupLocalName)
}

func testAccChannelWithOneRule(name, description, versionRange, actionName string) string {
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "test-project-group" {
			name = "%s"
		}

		resource "octopusdeploy_project" "test-project" {
			allow_deployments_to_no_targets = true
			lifecycle_id = "Lifecycles-1"
			name = "%s"
			project_group_id = octopusdeploy_project_group.test-project-group.id
		}

		resource "octopusdeploy_deployment_process" "deploy_step_template" {
			project_id = octopusdeploy_project.test-project.id

			step {
				name = "step-1"
				target_roles = ["Webserver",]

				deploy_package_action {
					features = [
						"Octopus.Features.ConfigurationTransforms",
						"Octopus.Features.ConfigurationVariables",
						"Octopus.Features.CustomDirectory",
						"Octopus.Features.CustomScripts",
						"Octopus.Features.IISWebSite"
					]
					name = "%s"

					primary_package {
						feed_id = "feeds-builtin"
						package_id = "MyPackage"
					}
				}
			}
		}

		resource "octopusdeploy_channel" "ch" {
		  depends_on  = ["octopusdeploy_deployment_process.deploy_step_template"]
		  description = "%s"
		  name = "%s"
		  project_id = octopusdeploy_project.test-project.id

		  rule {
		    actions = ["%s"]
		    version_range = "%s"
		  }
		}
		`,
		projectGroupName, projectName, actionName, description, name, actionName, versionRange,
	)
}

func testAccChannelWithTwoRules(name, description, versionRange1, actionName1, versionRange2, actionName2 string) string {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	projectTestOptions := NewProjectTestOptions(projectGroupLocalName, lifecycleLocalName)
	projectTestOptions.AllowDeploymentsToNoTargets = true

	return testAccProjectGroup(projectGroupLocalName, projectGroupName) + "\n" +
		testAccLifecycle(lifecycleLocalName, lifecycleName) + "\n" +
		testAccProjectWithOptions(projectTestOptions) + "\n" +
		fmt.Sprintf(`resource "octopusdeploy_deployment_process" "deploy_step_template" {
			project_id          = octopusdeploy_project.`+projectTestOptions.LocalName+`.id
			step {
				name            = "step-1"
				target_roles    = ["Webserver",]
				action {
					name 		= "%s"
					action_type = "Octopus.TentaclePackage"

					properties = {
						"Octopus.Action.Package.FeedId": "feeds-builtin"
						"Octopus.Action.Package.PackageId": "#{PackageName}"
					}

				}

				action {
					name 		= "%s"
					action_type = "Octopus.TentaclePackage"

					properties = {
						"Octopus.Action.Package.FeedId": "feeds-builtin"
						"Octopus.Action.Package.PackageId": "#{PackageName}"
					}

				}
			}
		}

		resource "octopusdeploy_channel" "ch" {
			description = "%s"
			name        = "%s"
			project_id  = octopusdeploy_project.`+projectTestOptions.LocalName+`.id

			rule {
				actions       = ["%s"]
				version_range = "%s"
			}

			rule {
				version_range = "%s"
				actions       = ["%s"]
			}

			depends_on = ["octopusdeploy_deployment_process.deploy_step_template"]
		}
		`,
			actionName1, actionName2, name, description, versionRange1, actionName1, versionRange2, actionName2,
		)
}

func testAccChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := existsHelperChannel(s, octoClient); err != nil {
			return err
		}
		return nil
	}
}

func existsHelperChannel(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_channel" {
			if _, err := client.Channels.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving channel %s", err)
			}
		}
	}
	return nil
}

func testAccChannelCheckDestroy(s *terraform.State) error {
	if err := destroyHelperChannel(s, octoClient); err != nil {
		return err
	}
	if err := testAccEnvironmentCheckDestroy(s); err != nil {
		return err
	}
	return nil
}

func destroyHelperChannel(s *terraform.State, client *client.Client) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_channel" {
			continue
		}

		if _, err := client.Channels.GetByID(rs.Primary.ID); err != nil {
			apiError := err.(*core.APIError)
			if apiError.StatusCode == http.StatusNotFound {
				continue
			}
			return fmt.Errorf("error retrieving channel %s", err)
		}
		return fmt.Errorf("channel still exists")
	}
	return nil
}

type ProjectTestOptions struct {
	AllowDeploymentsToNoTargets bool
	LifecycleLocalName          string
	LocalName                   string
	Name                        string
	ProjectGroupLocalName       string
}

func NewProjectTestOptions(projectGroupLocalName string, lifecycleLocalName string) *ProjectTestOptions {
	return &ProjectTestOptions{
		LifecycleLocalName:    lifecycleLocalName,
		LocalName:             acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Name:                  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ProjectGroupLocalName: projectGroupLocalName,
	}
}

func testAccProjectWithOptions(opt *ProjectTestOptions) string {

	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		allow_deployments_to_no_targets = %v
		lifecycle_id                    = octopusdeploy_lifecycle.%s.id
		name                            = "%s"
		project_group_id                = octopusdeploy_project_group.%s.id
	}`, opt.LocalName, opt.AllowDeploymentsToNoTargets, opt.LifecycleLocalName, opt.Name, opt.ProjectGroupLocalName)
}

func testAccProjectWithTemplate(localName string, name string, lifecycleLocalName string, projectGroupLocalName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		name             = "%s"
		project_group_id = octopusdeploy_project_group.%s.id

		template {
			name  = "project variable template name"
			label = "project variable template label"

			display_settings = {
				"Octopus.ControlType" = "Sensitive"
			}
		}
	}`, localName, lifecycleLocalName, name, projectGroupLocalName)

}
