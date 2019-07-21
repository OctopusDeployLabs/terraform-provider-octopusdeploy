package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployChannelBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_channel.ch"
	const channelName = "Funky Channel"
	const channelDescription = "this is Funky"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		//CheckDestroy: TODO,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelBasic(channelName, channelDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployChannelExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", channelName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", channelDescription),
				),
			},
		},
	})
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
			if _, err := client.Channel.Get(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving channel %s", err)
			}
		}
	}
	return nil
}

func testAccChannelBasic(name, description string) string {
	return fmt.Sprintf(`

		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			name           	= "funky project"
			lifecycle_id	= "Lifecycles-1"
			project_group_id = "${octopusdeploy_project_group.foo.id}" 	
			allow_deployments_to_no_targets = "True"
		}
		
		resource "octopusdeploy_channel" "ch" {
			name           	= "%s"
			description    	= "%s"
			project_id		= "${octopusdeploy_project.foo.id}"
		  }
		`,
		name, description,
	)
}