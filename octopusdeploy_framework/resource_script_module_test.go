package octopusdeploy_framework

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/scriptmodules"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployScriptModuleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_script_module." + localName

	body := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	syntax := "Bash"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testScriptModuleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testScriptModule(localName, name, description, body, syntax),
				Check: resource.ComposeTestCheckFunc(
					testScriptModuleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "script.#", "1"),
					resource.TestCheckResourceAttr(prefix, "script.0.body", body),
					resource.TestCheckResourceAttr(prefix, "script.0.syntax", syntax),
				),
			},
		},
	})
}

func testScriptModule(localName string, name string, description string, body string, syntax string) string {
	return fmt.Sprintf(`resource "octopusdeploy_script_module" "%s" {
		description = "%s"
		name        = "%s"

		script {
			body   = "%s"
			syntax = "%s"
		}
	}`, localName, description, name, body, syntax)
}

func testScriptModuleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		scriptModuleID := rs.Primary.ID
		if scriptModule, err := octoClient.ScriptModules.GetByID(scriptModuleID); err == nil {
			if scriptModule != nil {
				return fmt.Errorf("script module (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testScriptModuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_script_module" {
				if _, err := octoClient.ScriptModules.GetByID(r.Primary.ID); err != nil {
					return fmt.Errorf("error retrieving script module %s", err)
				}
			}
		}
		return nil
	}
}

// TestScriptModuleResource verifies that a script module set can be reimported with the correct settings
func TestScriptModuleResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "23-scriptmodule", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "23a-scriptmoduleds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := variables.LibraryVariablesQuery{
		PartialName: "Test2",
		Skip:        0,
		Take:        1,
	}

	resources, err := scriptmodules.Get(client, newSpaceId, query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a library variable set called \"Test2\"")
	}
	resource := resources.Items[0]

	if resource.Description != "Test script module" {
		t.Fatal("The library variable set must be have a description of \"Test script module\" (was \"" + resource.Description + "\")")
	}

	if resource.Syntax != "PowerShell" {
		t.Fatal("The script module must have a syntax of \"PowerShell\" (was \"" + resource.Syntax + "\")")
	}

	if resource.ScriptBody != "echo \"hi\"" {
		t.Fatal("The script module must have a script body of \"echo \"hi\"\" (was \"" + resource.ScriptBody + "\")")
	}

	variables, err := client.Variables.GetAll(resource.ID)

	if len(variables.Variables) != 2 {
		t.Fatal("The library variable set must have two associated variables")
	}

	foundScript := false
	foundLanguage := false
	for _, u := range variables.Variables {
		if u.Name == "Octopus.Script.Module[Test2]" {
			foundScript = true

			if u.Type != "String" {
				t.Fatal("The library variable set variable must have a type of \"String\"")
			}

			if u.Value != "echo \"hi\"" {
				t.Fatal("The library variable set variable must have a value of \"\"echo \\\"hi\\\"\"\"")
			}

			if u.IsSensitive {
				t.Fatal("The library variable set variable must not be sensitive")
			}

			if !u.IsEditable {
				t.Fatal("The library variable set variable must be editable")
			}
		}

		if u.Name == "Octopus.Script.Module.Language[Test2]" {
			foundLanguage = true

			if u.Type != "String" {
				t.Fatal("The library variable set variable must have a type of \"String\"")
			}

			if u.Value != "PowerShell" {
				t.Fatal("The library variable set variable must have a value of \"PowerShell\"")
			}

			if u.IsSensitive {
				t.Fatal("The library variable set variable must not be sensitive")
			}

			if !u.IsEditable {
				t.Fatal("The library variable set variable must be editable")
			}
		}
	}

	if !foundLanguage || !foundScript {
		t.Fatal("Script module must create two variables for script and language")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "23a-scriptmoduleds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
