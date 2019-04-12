package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployDeploymentTargetTriggerAddDelete(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project_deployment_target_trigger.foo"
	const deployTargetTriggerName = "Funky Monkey Trigger"
	projectName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDeploymentTargetTriggerResource(t, deployTargetTriggerName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectTriggerExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", deployTargetTriggerName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "should_redeploy", "true"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_groups.0", "Machine"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_categories.0", "MachineCleanupFailed"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDeploymentTargetTriggerUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project_deployment_target_trigger.foo"
	const deployTargetTriggerName = "Funky Monkey Trigger"
	projectName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDeploymentTargetTriggerResource(t, deployTargetTriggerName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectTriggerExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_groups.0", "Machine"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_categories.0", "MachineCleanupFailed"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "should_redeploy", "true"),
				),
			},
			{
				Config: testAccProjectDeploymentTargetTriggerResourceUpdated(t, deployTargetTriggerName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectTriggerExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_groups.0", "Machine"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_groups.1", "MachineCritical"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "event_categories.0", "MachineHealthy"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "should_redeploy", "false"),
				),
			},
		},
	})
}

func testAccProjectDeploymentTargetTriggerResource(t *testing.T, triggerName, projectName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			lifecycle_id          = "Lifecycles-1"
			name                  = "%s"
			project_group_id      = "${octopusdeploy_project_group.foo.id}"
	  	}

		resource "octopusdeploy_project_deployment_target_trigger" "foo" {
			name             = "%s"
			project_id       = "${octopusdeploy_project.foo.id}"
			event_groups     = ["Machine"]
			event_categories = ["MachineCleanupFailed"]
			should_redeploy  = true

			roles = [
			"FooRoles"
			]
		}
		`,
		projectName, triggerName,
	)
}

func testAccProjectDeploymentTargetTriggerResourceUpdated(t *testing.T, triggerName, projectName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_project_group" "foo" {
			name = "Integration Test Project Group"
		}

		resource "octopusdeploy_project" "foo" {
			lifecycle_id          = "Lifecycles-1"
			name                  = "%s"
			project_group_id      = "${octopusdeploy_project_group.foo.id}"
	  	}

		resource "octopusdeploy_project_deployment_target_trigger" "foo" {
			name             = "%s"
			project_id       = "${octopusdeploy_project.foo.id}"
			event_groups     = ["Machine", "MachineCritical"]
			event_categories = ["MachineHealthy"]
			should_redeploy  = false

			roles = [
			"FooRoles"
			]
		}
		`,
		projectName, triggerName,
	)
}

// testAccProjectDeploymentTriggerExists checks if a ProjectTrigger Exists
func testAccProjectTriggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)

		if _, err := client.ProjectTrigger.Get(rs.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving project trigger %s", err)
		}

		return nil
	}
}
