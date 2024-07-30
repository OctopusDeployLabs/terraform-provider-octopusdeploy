package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTenantProjectResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.22.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: tenantProjectConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   tenantProjectConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedTenantProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testFeedUpdated(t),
				),
			},
		},
	})
}

const tenantProjectConfig = `
	resource "octopusdeploy_tenant" "tenant1" {
		name        = "tenant 1"
	}

	resource "octopusdeploy_project" "project1" {
		name = "project 1"	
		lifecycle_id = "Lifecycles-1"
		project_group_id = "ProjectGroups-1"
	}

	resource "octopusdeploy_environment" "environment1" {
		name = "environment 1"
	}

	resource "octopusdeploy_environment" "environment2" {
		name = "environment 2"
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id = octopusdeploy_tenant.tenant1.id
		project_id   = octopusdeploy_project.project1.id
		environment_ids = [octopusdeploy_environment.environment1.id, octopusdeploy_environment.environment2.id]
    }`

const updatedTenantProjectConfig = `
	resource "octopusdeploy_tenant" "tenant1" {
		name        = "tenant 1"
	}

	resource "octopusdeploy_project" "project1" {
		name = "project 1"	
		lifecycle_id = "Lifecycles-1"
		project_group_id = "ProjectGroups-1"
	}

	resource "octopusdeploy_environment" "environment1" {
		name = "environment 1"
	}

	resource "octopusdeploy_environment" "environment2" {
		name = "environment 2"
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id = octopusdeploy_tenant.tenant1.id
		project_id   = octopusdeploy_project.project1.id
		environment_ids = [octopusdeploy_environment.environment1.id]
    }`

func testTenantProjectDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_maven_feed" {
			continue
		}

		feed, err := octoClient.Tenants.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("tenant (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testTenantProjectUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tenantId := s.RootModule().Resources["octopusdeploy_tenant_project"+".tenant_project_migration"].Primary.ID
		tenant, err := octoClient.Tenants.GetByID(tenantId)
		if err != nil {
			return fmt.Errorf("failed to retrieve tenant by ID: %s", err)
		}

		assert.Equal(t, "Tenants-1", tenant.ID, "Tenant ID did not match expected value")
		assert.Equal(t, len(tenant.ProjectEnvironments["Projects-1"]), 1, "environments collection should only have 1 entry")

		return nil
	}
}
