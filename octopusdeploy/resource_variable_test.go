package octopusdeploy

import (
	"fmt"
	"os"
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
		Providers:    testAccProviders,
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

		client := testAccProvider.Meta().(*client.Client)
		if _, err := client.Variables.GetByID(ownerID, variableID); err != nil {
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

	client := testAccProvider.Meta().(*client.Client)
	variable, err := client.Variables.GetByID(ownerID, variableID)
	if err == nil {
		if variable != nil {
			return fmt.Errorf("variable (%s) still exists", variableID)
		}
	}

	return nil
}
