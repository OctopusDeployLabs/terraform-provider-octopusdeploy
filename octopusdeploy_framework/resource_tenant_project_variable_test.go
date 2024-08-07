package octopusdeploy_framework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTenantProjectVariableBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryEnvironmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryEnvironmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryEnvironmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryEnvironmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	primaryResourceName := "octopusdeploy_tenant_project_variable." + primaryLocalName
	secondaryResourceName := "octopusdeploy_tenant_project_variable." + secondaryLocalName

	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantProjectVariableCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(primaryResourceName),
					testTenantProjectVariableExists(secondaryResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "value", primaryValue),
					resource.TestCheckResourceAttr(secondaryResourceName, "value", secondaryValue),
				),
				Config: testAccTenantProjectVariable(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, primaryEnvironmentLocalName, primaryEnvironmentName, secondaryEnvironmentLocalName, secondaryEnvironmentName, tenantLocalName, tenantName, tenantDescription, primaryLocalName, primaryValue, secondaryLocalName, secondaryValue),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(primaryResourceName),
					testTenantProjectVariableExists(secondaryResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "value", primaryValue),
					resource.TestCheckResourceAttr(secondaryResourceName, "value", newValue),
				),
				Config: testAccTenantProjectVariable(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, primaryEnvironmentLocalName, primaryEnvironmentName, secondaryEnvironmentLocalName, secondaryEnvironmentName, tenantLocalName, tenantName, tenantDescription, primaryLocalName, primaryValue, secondaryLocalName, newValue),
			},
		},
	})
}

func testAccTenantProjectVariable(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, primaryEnvironmentLocalName string, primaryEnvironmentName string, secondaryEnvironmentLocalName string, secondaryEnvironmentName string, tenantLocalName string, tenantName string, tenantDescription string, primaryLocalName string, primaryValue string, secondaryLocalName string, secondaryValue string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	var tfconfig = testAccLifecycle(lifecycleLocalName, lifecycleName) + "\n" +
		testAccProjectGroup(projectGroupLocalName, projectGroupName) + "\n" +
		testAccProjectWithTemplate(projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName) + "\n" +
		testAccEnvironment(primaryEnvironmentLocalName, primaryEnvironmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccEnvironment(secondaryEnvironmentLocalName, secondaryEnvironmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccTenantWithProjectEnvironment(tenantLocalName, tenantName, projectLocalName, primaryEnvironmentLocalName, secondaryEnvironmentLocalName) + "\n" +
		testTenantProjectVariable(primaryLocalName, primaryEnvironmentLocalName, projectLocalName, tenantLocalName, projectLocalName, primaryValue) + "\n" +
		testTenantProjectVariable(secondaryLocalName, secondaryEnvironmentLocalName, projectLocalName, tenantLocalName, projectLocalName, secondaryValue)
	return tfconfig
}

func testAccProjectWithTemplate(localName string, name string, lifecycleLocalName string, projectGroupLocalName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		name             = "%s"
		project_group_id = octopusdeploy_project_group.%s.id

		template {
			name  = "project variable template name"
			label = "project variable template label"

			display_settings = {
				"Octopus.ControlType" = "Sensitive"
			}
		}
	}`, localName, lifecycleLocalName, name, projectGroupLocalName)
}

func testAccProjectGroup(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccTenantWithProjectEnvironment(localName string, name string, projectLocalName string, primaryEnvironmentLocalName string, secondaryEnvironmentLocalName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id = octopusdeploy_tenant.%s.id
		project_id   = "${octopusdeploy_project.%s.id}"
		environment_ids = [octopusdeploy_environment.%s.id, octopusdeploy_environment.%s.id]
	}`, localName, name, localName, projectLocalName, primaryEnvironmentLocalName, secondaryEnvironmentLocalName)
}

func testTenantProjectVariable(localName string, environmentLocalName string, projectLocalName string, tenantLocalName string, templateLocalName string, value string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant_project_variable" "%s" {
		environment_id = octopusdeploy_environment.%s.id
		project_id     = octopusdeploy_project.%s.id
		tenant_id      = octopusdeploy_tenant.%s.id
		template_id    = octopusdeploy_project.%s.template[0].id
		value          = "%s"
 		depends_on     = [
            octopusdeploy_project.%s,
            octopusdeploy_tenant_project.project_environment
        ]
	}`, localName, environmentLocalName, projectLocalName, tenantLocalName, templateLocalName, value, projectLocalName)
}

func testTenantProjectVariableExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var environmentID string
		var projectID string
		var templateID string
		var tenantID string

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_tenant_project_variable" {
				environmentID = r.Primary.Attributes["environment_id"]
				projectID = r.Primary.Attributes["project_id"]
				templateID = r.Primary.Attributes["template_id"]
				tenantID = r.Primary.Attributes["tenant_id"]
			}
		}

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
			if _, ok := projectVariable.Variables[environmentID]; ok {
				if _, ok := projectVariable.Variables[environmentID][templateID]; ok {
					return nil
				}
			}
		}

		return fmt.Errorf("tenant project variable not found")
	}
}

func testAccTenantProjectVariableCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_project_variable" {
			continue
		}

		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 4 {
			return fmt.Errorf("octopusdeploy_tenant_project_variable import must be in the form of TenantID:ProjectID:EnvironmentID:TemplateID (e.g. Tenants-123:Projects-456:Environments-789:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		projectID := importStrings[1]
		environmentID := importStrings[2]
		templateID := importStrings[3]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return nil
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return nil
		}

		if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
			if _, ok := projectVariable.Variables[environmentID]; ok {
				if _, ok := projectVariable.Variables[environmentID][templateID]; ok {
					return fmt.Errorf("tenant project variable (%s) still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}
