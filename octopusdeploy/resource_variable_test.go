package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_variable." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isSensitive := false
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	variableType := "String"

	accountVariableType := "GoogleCloudAccount"
	accountValue := "octopusdeploy_gcp_account." + localName + ".id"

	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	// TODO: replace with client reference
	spaceID := os.Getenv("OCTOPUS_SPACE")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testVariableDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVariableExists(),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "owner_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", value),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "6"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(spaceID, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, value, variableType),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVariableExists(),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "owner_id"),
					resource.TestCheckResourceAttr(prefix, "type", variableType),
					resource.TestCheckResourceAttr(prefix, "value", newValue),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "6"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(spaceID, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, newValue, variableType),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVariableExists(),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "owner_id"),
					resource.TestCheckResourceAttr(prefix, "type", accountVariableType),
					resource.TestCheckResourceAttr(prefix, "value", accountValue),
					resource.TestCheckResourceAttr(prefix, "scope.#", "1"),
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "6"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: fmt.Sprintf(`%s

%s`,
					testGcpAccount(localName, name, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)),
					testVariableBasic(spaceID, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, accountValue, accountVariableType)),
			},
		},
	})
}

func testVariableBasic(spaceID string, environmentLocalName string,
	environmentName string,
	lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, channelLocalName string, channelName string, localName string, name string, description string, isSensitive bool, value string, variableType string) string {
	return fmt.Sprintf(`%s

		%s

		%s

		%s

		%s

		resource "octopusdeploy_variable" "%s" {
		  description  = "%s"
		  is_sensitive = "%v"
		  name         = "%s"
		  owner_id     = octopusdeploy_project.%s.id
		  type         = "%s"
		  value        = "%s"

		  scope {
			channels     = [octopusdeploy_channel.%s.id]
		    environments = [octopusdeploy_environment.%s.id]
			tenant_tags  = []
		  }
		}`,
		createEnvironment(environmentLocalName, environmentName),
		createLifecycle(lifecycleLocalName, lifecycleName),
		createProjectGroup(projectGroupLocalName, projectGroupName),
		createProject(spaceID, projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName),
		createChannel(channelLocalName, channelName, projectLocalName),
		localName,
		description,
		isSensitive,
		name,
		projectLocalName,
		variableType,
		value,
		channelLocalName,
		environmentLocalName,
	)
}

func createEnvironment(localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}`, localName, name)
}

func createLifecycle(localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}`, localName, name)
}

func createProjectGroup(localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func createProject(spaceID string, localName, name, lifecycleLocal, projectGroupLocal string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		name             = "%s"
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		project_group_id = octopusdeploy_project_group.%s.id
		space_id 		 = "%s"
	}`, localName, name, lifecycleLocal, projectGroupLocal, spaceID)
}

func createChannel(localName, name, projectLocal string) string {
	return fmt.Sprintf(`resource "octopusdeploy_channel" "%s" {
		name = "%s"
		project_id = octopusdeploy_project.%s.id
	}`, localName, name, projectLocal)
}

func testAccCheckVariableExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var ownerID string
		var variableID string

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_project" {
				ownerID = r.Primary.ID
			}

			if r.Type == "octopusdeploy_variable" {
				ownerID = r.Primary.Attributes["owner_id"]
				variableID = r.Primary.ID
			}
		}

		if _, err := octoClient.Variables.GetByID(ownerID, variableID); err != nil {
			return fmt.Errorf("error retrieving variable %s", err)
		}

		return nil
	}
}

func testVariableDestroy(s *terraform.State) error {
	var ownerID string
	var variableID string

	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_project" {
			ownerID = r.Primary.ID
		}

		if r.Type == "octopusdeploy_variable" {
			variableID = r.Primary.ID
		}
	}

	variable, err := octoClient.Variables.GetByID(ownerID, variableID)
	if err == nil {
		if variable != nil {
			return fmt.Errorf("variable (%s) still exists", variableID)
		}
	}

	return nil
}

// TestVariableSetResource verifies that a variable set can be reimported with the correct settings
func TestVariableSetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "18-variableset", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "18a-variablesetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := variables.LibraryVariablesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.LibraryVariableSets.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a library variable set called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test variable set" {
			t.Fatal("The library variable set must be have a description of \"Test variable set\" (was \"" + resource.Description + "\")")
		}

		variableSet, err := client.Variables.GetAll(resource.ID)

		if len(variableSet.Variables) != 1 {
			t.Fatal("The library variable set must have one associated variable")
		}

		if variableSet.Variables[0].Name != "Test.Variable" {
			t.Fatal("The library variable set variable must have a name of \"Test.Variable\"")
		}

		if variableSet.Variables[0].Type != "String" {
			t.Fatal("The library variable set variable must have a type of \"String\"")
		}

		if variableSet.Variables[0].Description != "Test variable" {
			t.Fatal("The library variable set variable must have a description of \"Test variable\"")
		}

		if variableSet.Variables[0].Value != "test" {
			t.Fatal("The library variable set variable must have a value of \"test\"")
		}

		if variableSet.Variables[0].IsSensitive {
			t.Fatal("The library variable set variable must not be sensitive")
		}

		if !variableSet.Variables[0].IsEditable {
			t.Fatal("The library variable set variable must be editable")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "18a-variablesetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

func TestVariableResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "49-variables", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		project, err := client.Projects.GetByName("Test")
		variableSet, err := client.Variables.GetAll(project.ID)

		if err != nil {
			return err
		}

		if len(variableSet.Variables) != 7 {
			t.Fatalf("Expected 7 variables to be created.")
		}

		for _, variable := range variableSet.Variables {
			switch variable.Name {
			case "UnscopedVariable":
				if !variable.Scope.IsEmpty() {
					t.Fatalf("Expected UnscopedVariable to have no scope values.")
				}
			case "ActionScopedVariable":
				if len(variable.Scope.Actions) == 0 {
					t.Fatalf("Expected ActionScopedVariable to have action scope.")
				}
			case "ChannelScopedVariable":
				if len(variable.Scope.Channels) == 0 {
					t.Fatalf("Expected ChannelScopedVariable to have channel scope.")
				}
			case "EnvironmentScopedVariable":
				if len(variable.Scope.Environments) == 0 {
					t.Fatalf("Expected EnvironmentScopedVariable to have environment scope.")
				}
			case "MachineScopedVariable":
				if len(variable.Scope.Machines) == 0 {
					t.Fatalf("Expected MachineScopedVariable to have machine scope.")
				}
			case "ProcessScopedVariable":
				if len(variable.Scope.ProcessOwners) == 0 {
					t.Fatalf("Expected ProcessScopedVariable to have process scope.")
				}
			case "RoleScopedVariable":
				if len(variable.Scope.Roles) == 0 {
					t.Fatalf("Expected RoleScopedVariable to have role scope.")
				}
			}
		}

		return nil
	})
}
