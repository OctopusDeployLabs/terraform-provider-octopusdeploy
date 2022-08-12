package octopusdeploy

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
// 	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// func TestAccDeploymentTargetTriggerAddDelete(t *testing.T) {
// 	var projectTrigger triggers.ProjectTrigger
// 	lifecycleLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	lifecycleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectGroupLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectGroupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	triggerLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	triggerName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

// 	name := "octopusdeploy_project_deployment_target_trigger." + triggerLocalName

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		CheckDestroy: resource.ComposeTestCheckFunc(
// 			testAccLifecycleCheckDestroy,
// 			testAccProjectGroupCheckDestroy,
// 			testAccProjectCheckDestroy,
// 			testAccProjectDeploymentTriggerCheckDestroy,
// 		),
// 		Steps: []resource.TestStep{
// 			{
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccProjectTriggerExists(name, &projectTrigger),
// 					resource.TestCheckResourceAttr(name, "name", triggerName),
// 					resource.TestCheckResourceAttr(name, "should_redeploy", "true"),
// 					resource.TestCheckResourceAttr(name, "event_groups.0", "Machine"),
// 					resource.TestCheckResourceAttr(name, "event_categories.0", "MachineCleanupFailed"),
// 				),
// 				Config: testAccProjectDeploymentTargetTriggerResource(t, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectLocalName, projectGroupName, projectName, triggerLocalName, triggerName),
// 			},
// 		},
// 	})
// }

// func TestAccDeploymentTargetTriggerUpdate(t *testing.T) {
// 	var projectTrigger triggers.ProjectTrigger
// 	lifecycleLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	lifecycleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectGroupLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectGroupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	projectName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	triggerLocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	triggerName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

// 	name := "octopusdeploy_project_deployment_target_trigger." + triggerLocalName

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		CheckDestroy: resource.ComposeTestCheckFunc(
// 			testAccLifecycleCheckDestroy,
// 			testAccProjectGroupCheckDestroy,
// 			testAccProjectCheckDestroy,
// 			testAccProjectDeploymentTriggerCheckDestroy,
// 		),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccProjectDeploymentTargetTriggerResource(t, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, triggerLocalName, triggerName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccProjectTriggerExists(name, &projectTrigger),
// 					resource.TestCheckResourceAttr(name, "event_groups.0", "Machine"),
// 					resource.TestCheckResourceAttr(name, "event_categories.0", "MachineCleanupFailed"),
// 					resource.TestCheckResourceAttr(name, "should_redeploy", "true"),
// 				),
// 			},
// 			{
// 				Config: testAccProjectDeploymentTargetTriggerResourceUpdated(t, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, triggerLocalName, triggerName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccProjectTriggerExists(name, &projectTrigger),
// 					resource.TestCheckResourceAttr(name, "event_groups.0", "Machine"),
// 					resource.TestCheckResourceAttr(name, "event_groups.1", "MachineCritical"),
// 					resource.TestCheckResourceAttr(name, "event_categories.0", "MachineHealthy"),
// 					resource.TestCheckResourceAttr(name, "should_redeploy", "false"),
// 				),
// 			},
// 		},
// 	})
// }

// // func testAccProjectDeploymentTargetTriggerResource(t *testing.T, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, triggerLocalName string, triggerName string) string {
// // 	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
// // 		testAccProjectGroup(projectGroupLocalName, projectName)+"\n"+
// // 		testAccProject(lifecycleLocalName, projectGroupLocalName, projectLocalName, projectName)+"\n"+
// // 		`
// // 		resource octopusdeploy_project_deployment_target_trigger "%s" {
// // 			event_categories = ["MachineCleanupFailed"]
// // 			event_groups     = ["Machine"]
// // 			name             = "%s"
// // 			project_id       = "${octopusdeploy_project.%s.id}"
// // 			roles            = ["FooRoles"]
// // 			should_redeploy  = true
// // 		}
// // 		`,
// // 		triggerLocalName, triggerName, projectLocalName,
// // 	)
// // }

// // func testAccProjectDeploymentTargetTriggerResourceUpdated(t *testing.T, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, triggerLocalName string, triggerName string) string {
// // 	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
// // 		testAccProjectGroup(projectGroupLocalName, projectName)+"\n"+
// // 		testAccProject(lifecycleLocalName, projectGroupLocalName, projectLocalName, projectName)+"\n"+
// // 		`
// // 		resource octopusdeploy_project_deployment_target_trigger "%s" {
// // 			event_categories = ["MachineHealthy"]
// // 			event_groups     = ["Machine", "MachineCritical"]
// // 			name             = "%s"
// // 			project_id       = "${octopusdeploy_project.%s.id}"
// // 			roles            = ["FooRoles"]
// // 			should_redeploy  = false
// // 		}
// // 		`,
// // 		triggerLocalName, triggerName, projectLocalName,
// // 	)
// // }

// func testAccProjectTriggerExists(resourceName string, projectTrigger *triggers.ProjectTrigger) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[resourceName]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", resourceName)
// 		}

// 		client := testAccProvider.Meta().(*client.Client)
// 		resource, err := client.ProjectTriggers.GetByID(rs.Primary.ID)
// 		if err != nil {
// 			return err
// 		}

// 		*projectTrigger = *resource
// 		return nil
// 	}
// }

// func testAccProjectDeploymentTriggerCheckDestroy(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*client.Client)
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "octopusdeploy_project_deployment_target_trigger" {
// 			continue
// 		}

// 		if project, err := client.ProjectTriggers.GetByID(rs.Primary.ID); err == nil {
// 			return fmt.Errorf("project deployment trigger (%s) still exists", project.GetID())
// 		}
// 	}

// 	return nil
// }
