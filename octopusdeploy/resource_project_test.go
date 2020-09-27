package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOctopusDeployProjectBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const allowDeploymentsToNoTargets = constTrue
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constLifecycleID, lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constAllowDeploymentsToNoTargets, allowDeploymentsToNoTargets),
				),
			},
		},
	})
}

//nolint:govet
func TestAccOctopusDeployProjectWithUpdate(t *testing.T) {
	return

	const terraformNamePrefix = "octopusdeploy_project.foo"
	const projectName = "Funky Monkey"
	const lifeCycleID = "Lifecycles-1"
	const allowDeploymentsToNoTargets = constTrue
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployProjectDestroy,
		Steps: []resource.TestStep{
			// create project with no description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constLifecycleID, lifeCycleID),
				),
			},
			// create update it with a description + build steps
			// update again by remove its description
			{
				Config: testAccProjectBasic(projectName, lifeCycleID, allowDeploymentsToNoTargets),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, projectName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constLifecycleID, lifeCycleID),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constDescription, ""),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(
						terraformNamePrefix, "deployment_step.0.iis_website.0.step_name"),
				),
			},
		},
	})
}

func testAccProjectBasic(name, lifeCycleID, allowDeploymentsToNoTargets string) string {
	return fmt.Sprintf(`
		resource constOctopusDeployProjectGroup "foo" {
			name = "Integration Test Project Group"
		}

		resource constOctopusDeployProject "foo" {
			name           = "%s"
			lifecycle_id    = "%s"
			project_group_id = "${octopusdeploy_project_group.foo.id}"
			allow_deployments_to_no_targets = "%s"
		}
		`,
		name, lifeCycleID, allowDeploymentsToNoTargets,
	)
}

func testAccCheckOctopusDeployProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	if err := destroyProjectHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyProjectHelper(s *terraform.State, apiClient *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := apiClient.Projects.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving project %s", err)
		}
		return fmt.Errorf("Project still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopus_deploy_project" {
			if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("Received an error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
