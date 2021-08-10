package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
		},
	})
}

func TestAccProjectWithUpdate(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.iis_website.0.step_name"),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
		},
	})
}

func testAccProject(lifecycleLocalName string, projectGroupLocalName string, localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		lifecycle_id = "${octopusdeploy_lifecycle.%s.id}"
		name = "%s"
		project_group_id = "${octopusdeploy_project_group.%s.id}"
	}`, localName, lifecycleLocalName, name, projectGroupLocalName)
}

func testAccProjectBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroupBasic(projectGroupLocalName, projectGroupName)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = "${octopusdeploy_lifecycle.%s.id}"
			name             = "%s"
			project_group_id = "${octopusdeploy_project_group.%s.id}"

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

func testAccProjectCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project" {
			continue
		}

		if project, err := client.Projects.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project (%s) still exists", project.GetID())
		}
	}

	return nil
}

func testAccProjectCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func existsHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_project" {
			if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
