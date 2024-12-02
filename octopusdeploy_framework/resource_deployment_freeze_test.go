package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestNewDeploymentFreezeResource(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_freeze." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	start := fmt.Sprintf("%d-11-21T06:30:00+10:00", time.Now().Year()+1)
	end := fmt.Sprintf("%d-11-21T08:30:00+10:00", time.Now().Year()+1)
	updatedEnd := fmt.Sprintf("%d-11-21T08:30:00+10:00", time.Now().Year()+2)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testDeploymentFreezeCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", end)),
				Config: testDeploymentFreezeBasic(localName, name, start, end, spaceName, environmentName, projectName, projectGroupName, lifecycleName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"1"),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", updatedEnd)),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, environmentName, projectName, projectGroupName, lifecycleName),
			},
		},
	})
}

func testDeploymentFreezeBasic(localName string, freezeName string, start string, end string, spaceName string, environmentName string, projectName string, projectGroupName string, lifecycleName string) string {
	spaceLocalName := fmt.Sprintf("space_%s", localName)
	projectScopeLocalName := fmt.Sprintf("project_scope_%s", localName)
	projectLocalName := fmt.Sprintf("project_%s", localName)
	lifecycleLocalName := fmt.Sprintf("lifecycle_%s", localName)
	projectGroupLocalName := fmt.Sprintf("project_group_%s", localName)
	environmentLocalName := fmt.Sprintf("environment_%s", localName)

	projectScopes := fmt.Sprintf(`resource "octopusdeploy_deployment_freeze_project" "%s" {
		deploymentfreeze_id = octopusdeploy_deployment_freeze.%s.id
		project_id = octopusdeploy_project.%s.id
		environment_ids = [ %s ]
	}`, projectScopeLocalName, localName, projectLocalName, fmt.Sprintf("octopusdeploy_environment.%s.id", environmentLocalName))

	return fmt.Sprintf(`
	%s

	%s 

	%s

	%s

	%s

	resource "octopusdeploy_deployment_freeze" "%s" {
		name = "%s"
		start = "%s"
		end = "%s"
	}

	%s`,
		createSpace(spaceLocalName, spaceName),
		createEnvironment(spaceLocalName, environmentLocalName, environmentName),
		createLifecycle(spaceLocalName, lifecycleLocalName, lifecycleName),
		createProjectGroup(spaceLocalName, projectGroupLocalName, projectGroupName),
		createProject(spaceLocalName, projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName),
		localName, freezeName, start, end, projectScopes)
}

//func testDeploymentFreezeBasic(localName string, freezeName string, start string, end string, spaceName string, environments []string, projectName string, projectGroupName string, lifecycleName string) string {
//	spaceLocalName := fmt.Sprintf("space_%s", localName)
//	projectScopeLocalName := fmt.Sprintf("project_scope_%s", localName)
//	projectLocalName := fmt.Sprintf("project_%s", localName)
//	lifecycleLocalName := fmt.Sprintf("lifecycle_%s", localName)
//	projectGroupLocalName := fmt.Sprintf("project_group_%s", localName)
//	environmentScopes := make([]string, 0, len(environments))
//	environmentResources := ""
//	for i, environmentName := range environments {
//		environmentLocalName := fmt.Sprintf("environment_%d_%s", i, localName)
//		environmentResources += fmt.Sprintln(createEnvironment(spaceLocalName, environmentLocalName, environmentName))
//		environmentScopes = append(environmentScopes, fmt.Sprintf("resource.octopusdeploy_environment.%s.id", environmentLocalName))
//	}
//
//	projectScopes := fmt.Sprintf(`resource "octopusdeploy_deployment_freeze_project" "%s" {
//		deploymentfreeze_id = octopusdeploy_deployment_freeze.%s.id
//		project_id = octopusdeploy_project.%s.id
//		environment_ids = [ %s ]
//	}`, projectScopeLocalName, localName, projectLocalName, strings.Join(environmentScopes, ","))
//
//	return fmt.Sprintf(`
//	%s
//
//	%s
//
//	%s
//
//	%s
//
//	%s
//
//	resource "octopusdeploy_deployment_freeze" "%s" {
//		name = "%s"
//		start = "%s"
//		end = "%s"
//	}
//
//	%s`,
//		createSpace(spaceLocalName, spaceName),
//		environmentResources,
//		createLifecycle(spaceLocalName, lifecycleLocalName, lifecycleName),
//		createProjectGroup(spaceLocalName, projectGroupLocalName, projectGroupName),
//		createProject(spaceLocalName, projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName),
//		localName, freezeName, start, end, projectScopes)
//}

func testDeploymentFreezeExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		freezeId := s.RootModule().Resources[prefix].Primary.ID
		if _, err := deploymentfreezes.GetById(octoClient, freezeId); err != nil {
			return err
		}

		return nil
	}
}

func testDeploymentFreezeCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_deployment_freeze" {
			continue
		}

		feed, err := deploymentfreezes.GetById(octoClient, rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Deployment Freeze (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
