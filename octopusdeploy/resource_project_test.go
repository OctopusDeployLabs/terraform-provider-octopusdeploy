package octopusdeploy

import (
	"fmt"
	"os"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectBasic(t *testing.T) {
	lifecycleTestOptions := test.NewLifecycleTestOptions()
	projectGroupTestOptions := test.NewProjectGroupTestOptions()
	projectTestOptions := test.NewProjectTestOptions(lifecycleTestOptions, projectGroupTestOptions)
	projectTestOptions.Resource.IsDisabled = true

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
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "description", projectTestOptions.Resource.Description),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "lifecycle_id", projectTestOptions.Resource.LifecycleID),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "name", projectTestOptions.Resource.Name),
				),
				Config: test.GetConfiguration([]string{
					test.LifecycleConfiguration(lifecycleTestOptions),
					test.ProjectGroupConfiguration(projectGroupTestOptions),
					test.ProjectConfiguration(projectTestOptions),
				}),
			},
		},
	})
}

func TestAccProjectCaC(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project." + localName

	basePath := ".octopus/" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) // note: replace with password
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)  // note: replace with valid space ID
	url := "https://example.com"                                        // note: replace with git repository address
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) // note: replace with valid username

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
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.0.base_path", basePath),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.0.credentials.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.0.credentials.0.password", password),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.0.credentials.0.username", username),
					resource.TestCheckResourceAttr(resourceName, "git_persistence_settings.0.url", url),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "lifecycle_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_group_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "space_id", spaceID),
				),
				Config: testAccProjectCaC(spaceID, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description, basePath, url, password, username),
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

	// TODO: replace with client reference
	spaceID := os.Getenv("OCTOPUS_SPACE")

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
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(spaceID, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.iis_website.0.step_name"),
				),
				Config: testAccProjectBasic(spaceID, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
		},
	})
}

func testAccProjectBasic(spaceID string, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	projectGroup := test.NewProjectGroupTestOptions()

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		test.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id
			space_id         = "%s"

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
		}`, localName, description, lifecycleLocalName, name, projectGroupLocalName, spaceID)
}

func testAccProjectCaC(spaceID string, lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string, basePath string, url string, password string, username string) string {
	projectGroup := test.NewProjectGroupTestOptions()

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		test.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
	       description      = "%s"
		   lifecycle_id     = octopusdeploy_lifecycle.%s.id
		   name             = "%s"
		   project_group_id = octopusdeploy_project_group.%s.id
		   space_id         = "%s"

			git_persistence_settings {
				base_path = "%s"
				url       = "%s"

				credentials {
					password = "%s"
					username = "%s"
				}
			}
		}`, localName, description, lifecycleLocalName, name, projectGroupLocalName, spaceID, basePath, url, password, username)
}

func testAccProjectCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
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

func testAccProjectCheckExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

func existsHelper(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_project" {
			if _, err := client.Projects.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving project with ID %s: %s", r.Primary.ID, err)
			}
		}
	}
	return nil
}
