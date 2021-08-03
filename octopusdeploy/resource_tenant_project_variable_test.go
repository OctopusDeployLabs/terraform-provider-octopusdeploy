package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTenantProjectVariableBasic(t *testing.T) {
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

	resourceName := "octopusdeploy_tenant_project_variable." + tenantVariablesLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(resourceName),
				),
				Config: testAccTenantProjectVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, value),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(resourceName),
				),
				Config: testAccTenantProjectVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, newValue),
			},
		},
	})
}

func testAccTenantProjectVariableBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, environmentLocalName string, environmentName string, tenantLocalName string, tenantName string, tenantDescription string, localName string, value string) string {
	return fmt.Sprintf(testAccLifecycleBasic(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroupBasic(projectGroupLocalName, projectGroupName)+"\n"+
		testEnvironmentMinimum(environmentLocalName, environmentName)+"\n"+`
		resource "octopusdeploy_project" "%s" {
			lifecycle_id                   = octopusdeploy_lifecycle.%s.id
			name                           = "%s"
			project_group_id               = octopusdeploy_project_group.%s.id

			template {
				name  = "project variable template name"
				label = "project variable template label"
			}
		}

		resource "octopusdeploy_tenant" "%s" {
			name = "%s"

			project_environment {
				project_id   = octopusdeploy_project.%s.id
				environments = [octopusdeploy_environment.%s.id]
			}
		}

		resource "octopusdeploy_tenant_project_variable" "%s" {
			environment_id = octopusdeploy_environment.%s.id
			is_sensitive   = false
			project_id     = octopusdeploy_project.%s.id
			tenant_id      = octopusdeploy_tenant.%s.id
			template_id    = octopusdeploy_project.%s.template[0].id
			value          = "%s"
		}`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, projectLocalName, environmentLocalName, localName, environmentLocalName, projectLocalName, tenantLocalName, projectLocalName, value)
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

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		tenant, err := client.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := client.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		for _, v := range tenantVariables.ProjectVariables {
			if v.ProjectID == projectID {
				for k := range v.Variables {
					if k == environmentID {
						if _, ok := v.Variables[environmentID][templateID]; ok {
							return nil
						}
					}
				}
			}
		}

		return fmt.Errorf("tenant project variable not found")
	}
}
