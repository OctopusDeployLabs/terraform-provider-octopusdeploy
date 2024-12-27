package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceBuiltInTrigger(t *testing.T) {
	localName := acctest.RandStringFromCharSet(50, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("octopusdeploy_built_in_trigger.%s", localName)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testBuiltInTriggerCheckDestroy(s) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configTestAccBuiltInTrigger(localName, "Action One", "console.one"),
				Check:  testAssertBuiltInResourceAttributesAreSet(prefix),
			},
			{
				Config: configTestAccBuiltInTrigger(localName, "Action Two", "console.two"),
				Check:  testAssertBuiltInResourceAttributesAreSet(prefix),
			},
		},
	})
}

func testAssertBuiltInResourceAttributesAreSet(prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(prefix, "project_id"),
		resource.TestCheckResourceAttrSet(prefix, "channel_id"),
		resource.TestCheckResourceAttrSet(prefix, "release_creation_package.deployment_action"),
		resource.TestCheckResourceAttrSet(prefix, "release_creation_package.package_reference"),
		testBuiltInTriggerExists(prefix),
	)
}

func testBuiltInTriggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		triggerResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		projectId, ok := triggerResource.Primary.Attributes["project_id"]
		if !ok {
			return fmt.Errorf("no project_id found: %s", resourceName)
		}

		project, failure := projects.GetByID(octoClient, octoClient.GetSpaceID(), projectId)
		if failure != nil {
			return fmt.Errorf("failed to read project(%s) of built-in trigger: %s", projectId, failure.Error())
		}

		if project.AutoCreateRelease == false {
			return fmt.Errorf("expected project.AutoCreateRelease not to be false")
		}

		if project.ReleaseCreationStrategy == nil {
			return fmt.Errorf("expected project.ReleaseCreationStrategy not to be nil")
		}

		return nil
	}
}

func testBuiltInTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_built_in_trigger" {
			continue
		}

		projectId := rs.Primary.Attributes["project_id"]
		project, err := projects.GetByID(octoClient, octoClient.GetSpaceID(), projectId)
		if err == nil && project != nil {
			return fmt.Errorf("project of built-in trigger (%s) still exists", projectId)
		}
	}

	return nil
}

func configTestAccBuiltInTrigger(localName string, actionName string, packageReference string) string {
	return fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}
		
		resource "octopusdeploy_project_group" "test" {
		  name        = "Test"
		  description = "Test project group"
		}

		resource "octopusdeploy_project" "test" {
		  name                                 = "Test"
		  lifecycle_id                         = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  project_group_id                     = octopusdeploy_project_group.test.id
		  default_guided_failure_mode          = "EnvironmentDefault"
		  default_to_skip_if_already_installed = false
		  description                          = "Project with Built-In Trigger"
		  discrete_channel_release             = false
		  is_disabled                          = false
		  is_discrete_channel_release          = false
		  is_version_controlled                = false
		  tenanted_deployment_participation    = "Untenanted"
		  included_library_variable_sets       = []
		
		  connectivity_policy {
			allow_deployments_to_no_targets = false
			exclude_unhealthy_targets       = false
			skip_machine_behavior           = "SkipUnavailableMachines"
		  }
		}
		
		resource "octopusdeploy_channel" "test" {
			name = "Test Channel"
			project_id = octopusdeploy_project.test.id
			lifecycle_id = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		}
		
		data "octopusdeploy_feeds" "built_in_feed" {
		  feed_type    = "BuiltIn"
		  ids          = null
		  partial_name = ""
		  skip         = 0
		  take         = 1
		}
		
		resource "octopusdeploy_deployment_process" "test" {
		  project_id = octopusdeploy_project.test.id

		  step {
			condition           = "Success"
			name                = "Step One"
			package_requirement = "LetOctopusDecide"
			start_trigger       = "StartAfterPrevious"
			run_script_action {
			  condition                          = "Success"
			  is_disabled                        = false
			  is_required                        = true
			  name                               = "Action One"
			  script_body                        = <<-EOT
				  $ExtractedPath = $OctopusParameters["Octopus.Action.Package[console.one].ExtractedPath"]
				  Write-Host $ExtractedPath
				EOT
			  run_on_server                      = true
		
			  package {
				name                      = "console.one"
				package_id                = "console.one"
				feed_id                   = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
				acquisition_location      = "Server"
				extract_during_deployment = true
			  }
			}
		  }

		  step {
			condition           = "Success"
			name                = "Step Two"
			package_requirement = "LetOctopusDecide"
			start_trigger       = "StartAfterPrevious"
			run_script_action {
			  condition                          = "Success"
			  is_disabled                        = false
			  is_required                        = true
			  name                               = "Action Two"
			  script_body                        = <<-EOT
				  $ExtractedPath = $OctopusParameters["Octopus.Action.Package[console.two].ExtractedPath"]
				  Write-Host $ExtractedPath
				EOT
			  run_on_server                      = true
		
			  package {
				name                      = "console.two"
				package_id                = "console.two"
				feed_id                   = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
				acquisition_location      = "Server"
				extract_during_deployment = true
			  }
			}
		  }
		}
		
		resource "octopusdeploy_built_in_trigger" "%s" {
		  project_id = octopusdeploy_project.test.id
		  channel_id = octopusdeploy_channel.test.id
		  
		  release_creation_package = {
			deployment_action = "%s"
			package_reference = "%s"
		  }
		
		  depends_on = [
			octopusdeploy_project.test,
			octopusdeploy_channel.test,
			octopusdeploy_deployment_process.test
		  ]
		}
		`,
		localName,
		actionName,
		packageReference,
	)
}
