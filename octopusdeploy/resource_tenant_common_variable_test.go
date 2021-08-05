package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTenantCommonVariableBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantVariablesLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_common_variable." + tenantVariablesLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
				Config: testAccTenantCommonVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, value),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
				),
				Config: testAccTenantCommonVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, newValue),
			},
		},
	})
}

func testAccTenantCommonVariableBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, environmentLocalName string, environmentName string, tenantLocalName string, tenantName string, tenantDescription string, localName string, value string) string {
	return fmt.Sprintf(testAccLifecycleBasic(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroupBasic(projectGroupLocalName, projectGroupName)+"\n"+
		testEnvironmentMinimum(environmentLocalName, environmentName)+"\n"+`
		resource "octopusdeploy_library_variable_set" "test-library-variable-set" {
			name = "test"

			template {
				default_value = "Default Value???"
				help_text     = "This is the help text"
				label         = "Test Label"
				name          = "Test Template"

				display_settings = {
					"Octopus.ControlType" = "Sensitive"
				}
			}
		}

		resource "octopusdeploy_project" "%s" {
			included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set.id]
			lifecycle_id                   = octopusdeploy_lifecycle.%s.id
			name                           = "%s"
			project_group_id               = octopusdeploy_project_group.%s.id
		}

		resource "octopusdeploy_tenant" "%s" {
			name = "%s"

			project_environment {
				project_id   = octopusdeploy_project.%s.id
				environments = [octopusdeploy_environment.%s.id]
			}
		}

		resource "octopusdeploy_tenant_common_variable" "%s" {
			library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set.id
			template_id             = octopusdeploy_library_variable_set.test-library-variable-set.template[0].id
			tenant_id               = octopusdeploy_tenant.%s.id
			value                   = "%s"
		}`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, projectLocalName, environmentLocalName, localName, tenantLocalName, value)
}

func testTenantCommonVariableExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var tenantID string
		var libraryVariableSetID string
		var templateID string

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_tenant_common_variable" {
				libraryVariableSetID = r.Primary.Attributes["library_variable_set_id"]
				templateID = r.Primary.Attributes["template_id"]
				tenantID = r.Primary.Attributes["tenant_id"]
			}
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		tenant, err := client.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := client.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		for _, v := range tenantVariables.LibraryVariables {
			if v.LibraryVariableSetID == libraryVariableSetID {
				if _, ok := v.Variables[templateID]; ok {
					return nil
				}
			}
		}

		return fmt.Errorf("tenant common variable not found")
	}
}
