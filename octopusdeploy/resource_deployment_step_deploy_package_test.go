package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOctopusDeployDeploymentStepDeployPackageBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_deployment_step_deploy_package.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentStepDeployPackageBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "project_id", "Project-000"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "step_name", "Run Verify Deploy Package"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "feed_id", "Feed-000"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "package", "cleanup.yolo"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "target_roles.0", "MyRole1"),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "target_roles.0", "MyRole2"),
					),
			},
		},
	})
}

func testAccDeploymentStepDeployPackageBasic() string {
	return `
		resource "octopusdeploy_deployment_step_deploy_package" "foo" {
			project_id 				= "Project-000"
			step_name         = "Run Verify Deploy Package"
			feed_id           = "Feed-000"
			package           = "cleanup.yolo"

			target_roles = [
				"MyRole1",
				"MyRole2",
			]
		}
	`
}
