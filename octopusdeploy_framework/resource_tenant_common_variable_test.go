package octopusdeploy_framework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
)

func TestAccTenantCommonVariableBasic(t *testing.T) {
	//SkipCI(t, "A managed resource \"octopusdeploy_project_group\" \"ewtxiwplhaenzmhpaqyx\" has\n        not been declared in the root module.")
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
		CheckDestroy:             testAccTenantCommonVariableCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
	projectGroup := internalTest.NewProjectGroupTestOptions()
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false
	projectGroup.LocalName = projectGroupLocalName

	var tfConfig = fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internalTest.ProjectGroupConfiguration(projectGroup)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
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

        resource "octopusdeploy_project" "%[1]s" {
            included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set.id]
            lifecycle_id                   = octopusdeploy_lifecycle.%[2]s.id
            name                           = "%[3]s"
            project_group_id               = octopusdeploy_project_group.%[4]s.id
            depends_on                     = [octopusdeploy_library_variable_set.test-library-variable-set]
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id]
            depends_on       = [octopusdeploy_project.%[1]s, octopusdeploy_tenant.%[5]s, octopusdeploy_environment.%[7]s]
        }

        resource "octopusdeploy_tenant_common_variable" "%[8]s" {
            library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set.id
            template_id             = octopusdeploy_library_variable_set.test-library-variable-set.template[0].id
            tenant_id               = octopusdeploy_tenant.%[5]s.id
            value                   = "%[9]s"
            depends_on              = [octopusdeploy_library_variable_set.test-library-variable-set, octopusdeploy_tenant_project.project_environment]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, environmentLocalName, localName, value)
	return tfConfig
}

func testTenantCommonVariableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if len(rs.Primary.ID) == 0 {
			return fmt.Errorf("Library variable ID is not set")
		}

		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 3 {
			return fmt.Errorf("octopusdeploy_tenant_common_variable import must be in the form of TenantID:LibraryVariableSetID:VariableID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		libraryVariableSetID := importStrings[1]
		templateID := importStrings[2]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
			if _, ok := libraryVariable.Variables[templateID]; ok {
				return nil
			}
		}

		return fmt.Errorf("tenant common variable not found")
	}
}

func testAccTenantCommonVariableCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_common_variable" {
			continue
		}

		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 3 {
			return fmt.Errorf("octopusdeploy_tenant_common_variable import must be in the form of TenantID:LibraryVariableSetID:VariableID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		libraryVariableSetID := importStrings[1]
		templateID := importStrings[2]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return nil
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return nil
		}

		if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
			if _, ok := libraryVariable.Variables[templateID]; ok {
				return fmt.Errorf("tenant common variable (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
