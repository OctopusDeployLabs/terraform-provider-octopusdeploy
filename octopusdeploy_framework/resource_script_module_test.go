package octopusdeploy_framework

import (
	"fmt"
	"testing"

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
