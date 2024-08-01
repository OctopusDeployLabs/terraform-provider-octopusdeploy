package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	internaltest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestAccProjectBasic(t *testing.T) {
	lifecycleTestOptions := internaltest.NewLifecycleTestOptions()
	projectGroupTestOptions := internaltest.NewProjectGroupTestOptions()
	projectTestOptions := internaltest.NewProjectTestOptions(lifecycleTestOptions, projectGroupTestOptions)
	projectTestOptions.Resource.IsDisabled = true

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleTestOptions.Resource.Name),
					testProjectGroupExists(projectGroupTestOptions.QualifiedName),
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "description", projectTestOptions.Resource.Description),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "name", projectTestOptions.Resource.Name),
				),
				Config: internaltest.GetConfiguration([]string{
					internaltest.LifecycleConfiguration(lifecycleTestOptions),
					internaltest.ProjectGroupConfiguration(projectGroupTestOptions),
					internaltest.ProjectConfiguration(projectTestOptions),
				}),
			},
		},
	})
}

func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}

func testProjectGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
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
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
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
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
		},
	})
}

func testAccProjectBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

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

		  	versioning_strategy {
				template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
		  	}
		    connectivity_policy {
		    allow_deployments_to_no_targets = true
			skip_machine_behavior           = "None"
		  }

		  version_control_settings {
			default_branch = "foo"
			url            = "https://example.com/"
			username       = "bar"
		  }

		  versioning_strategy {
		    template = "alskdjaslkdj"
		  }
		}`, localName, description, lifecycleLocalName, name, projectGroupLocalName)
}

func testAccProjectCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project" {
			continue
		}

		if project, err := octoClient.Projects.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project (%s) still exists", project.GetID())
		}
	}

	return nil
}

func testAccProjectCheckExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_project" {
				if _, err := octoClient.Projects.GetByID(r.Primary.ID); err != nil {
					return fmt.Errorf("error retrieving project with ID %s: %s", r.Primary.ID, err)
				}
			}
		}
		return nil
	}
}

// TestProjectResource verifies that a project can be reimported with the correct settings
func TestProjectResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "19-project", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "19a-projectds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test project" {
		t.Fatal("The project must be have a description of \"Test project\" (was \"" + resource.Description + "\")")
	}

	if resource.AutoCreateRelease {
		t.Fatal("The project must not have auto release create enabled")
	}

	if resource.DefaultGuidedFailureMode != "EnvironmentDefault" {
		t.Fatal("The project must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + resource.DefaultGuidedFailureMode + "\")")
	}

	if resource.DefaultToSkipIfAlreadyInstalled {
		t.Fatal("The project must not have DefaultToSkipIfAlreadyInstalled enabled")
	}

	if resource.IsDisabled {
		t.Fatal("The project must not have IsDisabled enabled")
	}

	if resource.IsVersionControlled {
		t.Fatal("The project must not have IsVersionControlled enabled")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The project must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}

	if len(resource.IncludedLibraryVariableSets) != 0 {
		t.Fatal("The project must not have any library variable sets")
	}

	if resource.ConnectivityPolicy.AllowDeploymentsToNoTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.ExcludeUnhealthyTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
		t.Log("BUG: The project must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + resource.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "19a-projectds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}

func TestProjectInSpaceResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "19b-projectspace", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)

	spaces, err := spaces.GetAll(client)

	if err != nil {
		t.Fatal(err.Error())
	}
	idx := sort.Search(len(spaces), func(i int) bool { return spaces[i].Name == "Project Space Test" })
	space := spaces[idx]

	query := projects.ProjectsQuery{
		PartialName: "Test project in space",
		Skip:        0,
		Take:        1,
	}

	resources, err := projects.Get(client, space.ID, query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test project in space\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test project in space" {
		t.Fatal("The project must be have a description of \"Test project in space\" (was \"" + resource.Description + "\")")
	}

	if resource.AutoCreateRelease {
		t.Fatal("The project must not have auto release create enabled")
	}

	if resource.DefaultGuidedFailureMode != "EnvironmentDefault" {
		t.Fatal("The project must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + resource.DefaultGuidedFailureMode + "\")")
	}

	if resource.DefaultToSkipIfAlreadyInstalled {
		t.Fatal("The project must not have DefaultToSkipIfAlreadyInstalled enabled")
	}

	if resource.IsDisabled {
		t.Fatal("The project must not have IsDisabled enabled")
	}

	if resource.IsVersionControlled {
		t.Fatal("The project must not have IsVersionControlled enabled")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The project must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}

	if len(resource.IncludedLibraryVariableSets) != 0 {
		t.Fatal("The project must not have any library variable sets")
	}

	if resource.ConnectivityPolicy.AllowDeploymentsToNoTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.ExcludeUnhealthyTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
		t.Log("BUG: The project must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + resource.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
	}
}

// TestProjectWithGitUsernameExport verifies that a project can be reimported with the correct git settings
func TestProjectWithGitUsernameExport(t *testing.T) {
	if os.Getenv("GIT_CREDENTIAL") == "" {
		t.Fatal("The GIT_CREDENTIAL environment variable must be set")
	}

	if os.Getenv("GIT_USERNAME") == "" {
		t.Fatal("The GIT_USERNAME environment variable must be set")
	}

	testFramework := test.OctopusContainerTest{}
	_, err := testFramework.Act(t, octoContainer, "../terraform", "39-projectgitusername", []string{
		"-var=project_git_password=" + os.Getenv("GIT_CREDENTIAL"),
		"-var=project_git_username=" + os.Getenv("GIT_USERNAME"),
	})

	if err != nil {
		t.Fatal(err.Error())
	}

	// The client does not expose git credentials, so just test the import worked ok
}

// TestProjectWithDollarSignsExport verifies that a project can be reimported with terraform string interpolation
func TestProjectWithDollarSignsExport(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "40-escapedollar", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}

}

// TestProjectTerraformInlineScriptExport verifies that a project can be reimported with a terraform inline template step.
// See https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/478
func TestProjectTerraformInlineScriptExport(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "41-terraforminlinescript", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	deploymentProcess, err := client.DeploymentProcesses.GetByID(resource.DeploymentProcessID)

	if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Terraform.Template"].Value != "#test" {
		t.Fatalf("The inline Terraform template must be set to \"#test\"")
	}
}

// TestProjectTerraformPackageScriptExport verifies that a project can be reimported with a terraform package template step.
// See https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/478
func TestProjectTerraformPackageScriptExport(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "42-terraformpackagescript", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	deploymentProcess, err := client.DeploymentProcesses.GetByID(resource.DeploymentProcessID)

	if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Script.ScriptSource"].Value != "Package" {
		t.Fatalf("The Terraform template must be set deploy files from a package")
	}

	if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Terraform.TemplateDirectory"].Value != "blah" {
		t.Fatalf("The Terraform template directory must be set to \"blah\"")
	}
}

// TestProjectWithScriptActions verifies that a project with a plain script step can be applied and reapplied
func TestProjectWithScriptActions(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "45-projectwithscriptactions", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Do a second apply to catch the scenario documented at https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/509
	err = testFramework.TerraformApply(t, filepath.Join("../terraform", "45-projectwithscriptactions"), octoContainer.URI, newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test project" {
		t.Fatal("The project must be have a description of \"Test project\" (was \"" + resource.Description + "\")")
	}

	if resource.AutoCreateRelease {
		t.Fatal("The project must not have auto release create enabled")
	}

	if resource.DefaultGuidedFailureMode != "EnvironmentDefault" {
		t.Fatal("The project must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + resource.DefaultGuidedFailureMode + "\")")
	}

	if resource.DefaultToSkipIfAlreadyInstalled {
		t.Fatal("The project must not have DefaultToSkipIfAlreadyInstalled enabled")
	}

	if resource.IsDisabled {
		t.Fatal("The project must not have IsDisabled enabled")
	}

	if resource.IsVersionControlled {
		t.Fatal("The project must not have IsVersionControlled enabled")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The project must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}

	if len(resource.IncludedLibraryVariableSets) != 0 {
		t.Fatal("The project must not have any library variable sets")
	}

	if resource.ConnectivityPolicy.AllowDeploymentsToNoTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.ExcludeUnhealthyTargets {
		t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
	}

	if resource.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
		t.Log("BUG: The project must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + resource.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
	}

	deploymentProcess, err := client.DeploymentProcesses.GetByID(resource.DeploymentProcessID)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(deploymentProcess.Steps) != 1 {
		t.Fatal("The DeploymentProcess should have a single Deployment Step")
	}
	step := deploymentProcess.Steps[0]

	if len(step.Actions) != 3 {
		t.Fatal("The DeploymentProcess should have a three Deployment Actions")
	}

	if step.Actions[0].Name != "Pre Script Action" {
		t.Fatal("The first Deployment Action should be name \"Pre Script Action\" (was \"" + step.Actions[0].Name + "\")")
	}
	if step.Actions[1].Name != "Hello world (using PowerShell)" {
		t.Fatal("The second Deployment Action should be name \"Hello world (using PowerShell)\" (was \"" + step.Actions[1].Name + "\")")
	}
	if step.Actions[2].Name != "Post Script Action" {
		t.Fatal("The third Deployment Action should be name \"Post Script Action\" (was \"" + step.Actions[2].Name + "\")")
	}
}
