package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployChannelBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelBasic(name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", name),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", description),
					resource.TestCheckResourceAttrSet(terraformNamePrefix, "project_id"),
				),
			},
		},
	})
}

func TestAccOctopusDeployChannelBasicWithUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			// create baseline channel
			{
				Config: testAccChannelBasic(channelName, "this is funky"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", "this is funky"),
				),
			},
			// update channel with a new description
			{
				Config: testAccChannelBasic(channelName, "funky it is"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", "funky it is"),
				),
			},
		},
	})
}

func TestAccOctopusDeployChannelWithOneRule(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const channelDescription = "this is Funky"
	const actionName = "Funky Action"
	const versionRange = "1.0"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			{ // create channel with one rule
				Config: testAccChannelWithOneRule(channelName, channelDescription, versionRange, actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			{ // create baseline channel
				Config: testAccChannelWithOneRule(channelName, channelDescription, versionRange, actionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(terraformNamePrefix, "description", channelDescription),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.version_range", versionRange),
					resource.TestCheckResourceAttr(terraformNamePrefix, "rule.0.actions.0", actionName),
				),
			},
			{ // create updated channel with new values
				Config: testAccChannelWithOneRule(updatedChannelName, updatedChannelDescription, updatedVersionRange, updatedActionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
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
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const channelDescription = "this is Funky"
	const versionRange1 = "1.0"
	const actionName1 = "Funky Action"
	const versionRange2 = "2.0"
	const actionName2 = "Action-2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccChannelCheckDestroy,
		Steps: []resource.TestStep{
			{ // create channel with two rules
				Config: testAccChannelWithtwoRules(channelName, channelDescription, versionRange1, actionName1, versionRange2, actionName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
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

func testAccChannelBasic(name string, description string) string {
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testAccProjectBasic(projectLocalName, projectName, projectDescription)+"\n"+`		
		resource "octopusdeploy_channel" "ch" {
			description = "%s"
			name        = "%s"
			project_id  = "${octopusdeploy_project.%s.id}"
		}`, description, name, projectLocalName)
}

func testAccChannelWithOneRule(name, description, versionRange, actionName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			allow_deployments_to_no_targets = true
			lifecycle_id	                = "Lifecycles-1"
			name           	                = "funky project"
			project_group_id                = "${octopusdeploy_project_group.foo.id}" 	
		}

		resource "octopusdeploy_deployment_process" "deploy_step_template" {
			project_id          = "${octopusdeploy_project.foo.id}"
			step {
				name            = "step-1"
				target_roles    = ["Webserver",]
				action {
					action_type = "Octopus.TentaclePackage"
					name 		= "%s"
		
					property {
						key 	= "Octopus.Action.Package.FeedId"
						value 	= "feeds-builtin"
					}
		
					property {
						key 	= "Octopus.Action.Package.PackageId"
						value 	= "#{PackageName}"
					}
				}
			}
		}
		
		resource "octopusdeploy_project_channel" "ch" {
		  depends_on  = ["octopusdeploy_deployment_process.deploy_step_template"]
		  description = "%s"
		  name        = "%s"
		  project_id  = "${octopusdeploy_project.foo.id}"

		  rule {
		    actions       = ["%s"] 
		    version_range = "%s"
		  }
		}
		`,
		actionName, name, description, versionRange, actionName,
	)
}

func testAccChannelWithtwoRules(name, description, versionRange1, actionName1, versionRange2, actionName2 string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			name           	= "funky project"
			lifecycle_id	= "Lifecycles-1"
			project_group_id = "${octopusdeploy_project_group.foo.id}" 	
			allow_deployments_to_no_targets = true
		}

		resource "octopusdeploy_deployment_process" "deploy_step_template" {
			project_id          = "${octopusdeploy_project.foo.id}"
			step {
				name            = "step-1"
				target_roles    = ["Webserver",]
				action {
					name 		= "%s"
					action_type = "Octopus.TentaclePackage"
		
					property {
						key 	= "Octopus.Action.Package.FeedId"
						value 	= "feeds-builtin"
					}
		
					property {
						key 	= "Octopus.Action.Package.PackageId"
						value 	= "#{PackageName}"
					}

				}
				
				action {
					name 		= "%s"
					action_type = "Octopus.TentaclePackage"
		
					property {
						key 	= "Octopus.Action.Package.FeedId"
						value 	= "feeds-builtin"
					}
		
					property {
						key 	= "Octopus.Action.Package.PackageId"
						value 	= "#{PackageName}"
					}

				}
			}
		}
		
		resource "octopusdeploy_channel" "ch" {
			name           	= "%s"
			description    	= "%s"
			project_id		= "${octopusdeploy_project.foo.id}"
			rule {
				version_range 	= "%s"
				actions 		= ["%s"] 
			}
			
			rule {
				version_range 	= "%s"
				actions 		= ["%s"] 
			}

			depends_on = ["octopusdeploy_deployment_process.deploy_step_template"]
		}
		`,
		actionName1, actionName2, name, description, versionRange1, actionName1, versionRange2, actionName2,
	)
}

func testAccCheckOctopusDeployChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelperChannel(s, client); err != nil {
			return err
		}
		return nil
	}
}

func existsHelperChannel(s *terraform.State, client *octopusdeploy.Client) error {
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
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyHelperChannel(s, client); err != nil {
		return err
	}
	if err := testEnvironmentDestroy(s); err != nil {
		return err
	}
	return nil
}

func destroyHelperChannel(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Channels.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("error retrieving channel %s", err)
		}
		return fmt.Errorf("channel still exists")
	}
	return nil
}
