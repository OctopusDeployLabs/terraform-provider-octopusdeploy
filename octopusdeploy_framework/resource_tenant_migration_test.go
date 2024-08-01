package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTenantResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.21.1 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		CheckDestroy: testTenantProjectDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.22.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: tenantConfig(),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   tenantConfig(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedTenantResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testTenantResourceUpdated(t, name),
				),
			},
		},
	})
}

func tenantConfig() string {
	return fmt.Sprintf(`
	resource "octopusdeploy_tenant" "tenant1" {
		name        = "tenant test"
	}`)
}

func updatedTenantResourceConfig() string {
	return fmt.Sprintf(`
resource "octopusdeploy_tag_set" "tagset_tag1" {
  name        = "tag1"
  description = "Test tagset"
  sort_order  = 0
}

resource "octopusdeploy_tag" "tag_a" {
  name        = "a"
  color       = "#333333"
  description = "tag a"
  sort_order  = 2
  tag_set_id = octopusdeploy_tag_set.tagset_tag1.id
}

resource "octopusdeploy_tag" "tag_b" {
  name        = "b"
  color       = "#333333"
  description = "tag b"
  sort_order  = 3
  tag_set_id = octopusdeploy_tag_set.tagset_tag1.id
}

	resource "octopusdeploy_tenant" "tenant1" {
  		name        = "Updated tenant"
  		description = "Updated description"
  		tenant_tags = ["tag1/a", "tag1/b"]
  		space_id    = var.octopus_space_id
  		depends_on  = [octopusdeploy_tag.tag_a, octopusdeploy_tag.tag_b]
	}`)
}

func testTenantResourceUpdated(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tenantId := s.RootModule().Resources["octopusdeploy_tenant.tenant1"].Primary.ID
		projectId := s.RootModule().Resources["octopusdeploy_project.project1"].Primary.ID
		environmentId := s.RootModule().Resources["octopusdeploy_environment.environment1"].Primary.ID
		environment2Id := s.RootModule().Resources["octopusdeploy_environment.environment2"].Primary.ID
		tenant, err := octoClient.Tenants.GetByID(tenantId)
		if err != nil {
			return fmt.Errorf("failed to retrieve tenant by ID: %s", err)
		}

		assert.NotEmpty(t, tenant.ID, "Tenant ID did not match expected value")
		assert.Equal(t, tenant.Description, fmt.Sprintf("description %s", name))
		assert.Equal(t, len(tenant.ProjectEnvironments[projectId]), 2, "environments collection should only have 1 entry")
		assert.Equal(t, tenant.ProjectEnvironments[projectId][0], environmentId, "environments collection should contain id for environment1")
		assert.Equal(t, tenant.ProjectEnvironments[projectId][1], environment2Id, "environments collection should contain id for environment2")

		return nil
	}
}
