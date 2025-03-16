package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccOctopusDeployVariableBasic(t *testing.T) {
	internalTest.SkipCI(t, "After applying this test step, the refresh plan was not empty. scope.tenant_tags")
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

	spaceID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testVariableDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
					resource.TestCheckResourceAttr(prefix, "scope.0.%", "7"),
					resource.TestCheckResourceAttr(prefix, "scope.0.environments.#", "1"),
				),
				Config: testVariableBasic(spaceID, spaceName, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, value, variableType),
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
				Config: testVariableBasic(spaceID, spaceName, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, newValue, variableType),
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
					testVariableBasic(spaceID, spaceName, environmentLocalName, environmentName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channelLocalName, channelName, localName, name, description, isSensitive, accountValue, accountVariableType)),
			},
		},
	})
}

func testVariableBasic(spaceID string, spaceName string, environmentLocalName string,
	environmentName string,
	lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, channelLocalName string, channelName string, localName string, name string, description string, isSensitive bool, value string, variableType string) string {
	return fmt.Sprintf(`%s

		%s

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
          space_id = octopusdeploy_space.%s.id
		}`,
		createSpace(spaceID, spaceName),
		createEnvironment(spaceID, environmentLocalName, environmentName),
		createLifecycle(spaceID, lifecycleLocalName, lifecycleName),
		createProjectGroup(spaceID, projectGroupLocalName, projectGroupName),
		createProject(spaceID, projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName),
		createChannel(spaceID, channelLocalName, channelName, projectLocalName),
		localName,
		description,
		isSensitive,
		name,
		projectLocalName,
		variableType,
		value,
		channelLocalName,
		environmentLocalName,
		spaceID,
	)
}

func testGcpAccount(localName string, name string, jsonKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcp_account" "%s" {
		json_key = "%s"
		name = "%s"
	}`, localName, jsonKey, name)
}

func createSpace(localname string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_space" "%s" {
		name = "%s"
		space_managers_teams = ["teams-administrators"]
}`, localname, name)
}

func createEnvironment(spaceId, localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		name = "%s"
		space_id = octopusdeploy_space.%s.id
	}`, localName, name, spaceId)
}

func createLifecycle(spaceId, localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
		space_id = octopusdeploy_space.%s.id
	}`, localName, name, spaceId)
}

func createProjectGroup(spaceId, localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
		space_id = octopusdeploy_space.%s.id
	}`, localName, name, spaceId)
}

func createProject(spaceID string, localName, name, lifecycleLocal, projectGroupLocal string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		name             = "%s"
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		project_group_id = octopusdeploy_project_group.%s.id
		space_id 		 = octopusdeploy_space.%s.id
	}`, localName, name, lifecycleLocal, projectGroupLocal, spaceID)
}

func createChannel(spaceId string, localName, name, projectLocal string) string {
	return fmt.Sprintf(`resource "octopusdeploy_channel" "%s" {
		name = "%s"
		project_id = octopusdeploy_project.%s.id
		space_id = octopusdeploy_space.%s.id
	}`, localName, name, projectLocal, spaceId)
}

func testAccCheckVariableExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var ownerID string
		var variableID string
		var spaceID string

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_project" {
				ownerID = r.Primary.ID
				spaceID = r.Primary.Attributes["space_id"]
			}

			if r.Type == "octopusdeploy_variable" {
				ownerID = r.Primary.Attributes["owner_id"]
				spaceID = r.Primary.Attributes["space_id"]
				variableID = r.Primary.ID
			}
		}

		if _, err := variables.GetByID(octoClient, spaceID, ownerID, variableID); err != nil {
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
