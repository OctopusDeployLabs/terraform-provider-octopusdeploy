package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployScriptModuleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_script_module." + localName

	body := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	syntax := "Bash"

	resource.Test(t, resource.TestCase{
		CheckDestroy: testScriptModuleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		scriptModuleID := rs.Primary.ID
		if scriptModule, err := client.ScriptModules.GetByID(scriptModuleID); err == nil {
			if scriptModule != nil {
				return fmt.Errorf("script module (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testScriptModuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_script_module" {
				if _, err := client.ScriptModules.GetByID(r.Primary.ID); err != nil {
					return fmt.Errorf("error retrieving script module %s", err)
				}
			}
		}
		return nil
	}
}
