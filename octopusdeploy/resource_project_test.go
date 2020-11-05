package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployProjectBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployProject + "." + localName

	allowDeploymentsToNoTargets := constTrue
	description := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(localName, name, allowDeploymentsToNoTargets, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, constAllowDeploymentsToNoTargets, allowDeploymentsToNoTargets),
					resource.TestCheckResourceAttr(prefix, constDescription, description),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
		},
	})
}

func TestAccOctopusDeployProjectWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployProject + "." + localName

	allowDeploymentsToNoTargets := constTrue
	description := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testProjectDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectBasic(localName, name, allowDeploymentsToNoTargets, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, constAllowDeploymentsToNoTargets, allowDeploymentsToNoTargets),
					resource.TestCheckResourceAttr(prefix, constDescription, description),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
			},
			{
				Config: testAccProjectBasic(localName, name, allowDeploymentsToNoTargets, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployProjectExists(prefix),
					resource.TestCheckResourceAttr(prefix, constAllowDeploymentsToNoTargets, allowDeploymentsToNoTargets),
					resource.TestCheckResourceAttr(prefix, constDescription, description),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.iis_website.0.step_name"),
				),
			},
		},
	})
}

func testAccProjectBasic(localName string, name string, allowDeploymentsToNoTargets string, description string) string {
	lifecycleLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	lifecycleID := "${" + constOctopusDeployLifecycle + "." + lifecycleLocalName + ".id}"
	projectGroupID := "${" + constOctopusDeployProjectGroup + "." + projectGroupLocalName + ".id}"

	return fmt.Sprintf(testAccLifecycleBasic(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroupBasic(projectGroupLocalName, projectGroupName)+"\n"+
		`resource "%s" "%s" {
			allow_deployments_to_no_targets = "%s"
			description                     = "%s"
			lifecycle_id                    = "%s"
			name                            = "%s"
			project_group_id                = "%s"
		}`, constOctopusDeployProject, localName, allowDeploymentsToNoTargets, description, lifecycleID, name, projectGroupID)
}

func testProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		projectID := rs.Primary.ID
		project, err := client.Projects.GetByID(projectID)
		if err == nil {
			if project != nil {
				return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckOctopusDeployProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyProjectHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func destroyProjectHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving project %s", err)
		}
		return fmt.Errorf("Project still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == constOctopusDeployProject {
			if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("Received an error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
