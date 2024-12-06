package octopusdeploy_framework

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewDeploymentFreezeResource(t *testing.T) {
	if os.Getenv("TF_LOG") == "" {
		os.Setenv("TF_LOG", "INFO")
	}

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_freeze." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	start := fmt.Sprintf("%d-11-21T06:30:00+10:00", time.Now().Year()+1)
	end := fmt.Sprintf("%d-11-21T08:30:00+10:00", time.Now().Year()+1)
	updatedEnd := fmt.Sprintf("%d-11-21T08:30:00+10:00", time.Now().Year()+2)

	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName1 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testDeploymentFreezeCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", end)),
				Config: testDeploymentFreezeBasic(localName, name, start, end, spaceName, []string{environmentName1}, projectName, projectGroupName, lifecycleName, tenantName, false),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"1"),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", updatedEnd)),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, []string{environmentName1, environmentName2}, projectName, projectGroupName, lifecycleName, tenantName, false),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					testDeploymentFreezeTenantExists(fmt.Sprintf("octopusdeploy_deployment_freeze_tenant.tenant_%s", localName), t)),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, []string{environmentName1, environmentName2}, projectName, projectGroupName, lifecycleName, tenantName, true),
			},
		},
	})
}

func testDeploymentFreezeBasic(localName string, freezeName string, start string, end string, spaceName string, environments []string, projectName string, projectGroupName string, lifecycleName string, tenantName string, includeTenant bool) string {
	spaceLocalName := fmt.Sprintf("space_%s", localName)
	projectScopeLocalName := fmt.Sprintf("project_scope_%s", localName)
	projectLocalName := fmt.Sprintf("project_%s", localName)
	lifecycleLocalName := fmt.Sprintf("lifecycle_%s", localName)
	projectGroupLocalName := fmt.Sprintf("project_group_%s", localName)
	tenantLocalName := fmt.Sprintf("tenant_%s", localName)

	environmentScopes := make([]string, 0, len(environments))
	environmentResources := ""
	for i, environmentName := range environments {
		environmentLocalName := fmt.Sprintf("environment_%d_%s", i, localName)
		environmentResources += fmt.Sprintln(createEnvironment(spaceLocalName, environmentLocalName, environmentName))
		environmentScopes = append(environmentScopes, fmt.Sprintf("resource.octopusdeploy_environment.%s.id", environmentLocalName))
	}

	config := fmt.Sprintf(`
       
        # Space Configuration
        %s

        # Environment Configuration
        %s

        # Lifecycle Configuration
        %s

        # Project Group Configuration
        %s

        # Project Configuration
        %s

        resource "octopusdeploy_deployment_freeze" "%s" {
            name = "%s"
            start = "%s"
            end = "%s"
        }

        resource "octopusdeploy_deployment_freeze_project" "%s" {
            deploymentfreeze_id = octopusdeploy_deployment_freeze.%s.id
            project_id = octopusdeploy_project.%s.id
            environment_ids = [ %s ]
        }`,
		createSpace(spaceLocalName, spaceName),
		environmentResources,
		createLifecycle(spaceLocalName, lifecycleLocalName, lifecycleName),
		createProjectGroup(spaceLocalName, projectGroupLocalName, projectGroupName),
		createProject(spaceLocalName, projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName),
		localName, freezeName, start, end,
		projectScopeLocalName, localName, projectLocalName,
		strings.Join(environmentScopes, ","))

	if includeTenant {
		tenantConfig := fmt.Sprintf(`
            resource "octopusdeploy_tenant" "%[1]s" {
                name = "%[2]s"
                space_id = octopusdeploy_space.%[3]s.id
            }

            resource "octopusdeploy_tenant_project" "%[1]s_project" {
                tenant_id = octopusdeploy_tenant.%[1]s.id
                project_id = octopusdeploy_project.%[4]s.id
                environment_ids = [octopusdeploy_environment.environment_0_%[5]s.id]
                space_id = octopusdeploy_space.%[3]s.id
            }

            resource "octopusdeploy_deployment_freeze_tenant" "%[1]s" {
                deploymentfreeze_id = octopusdeploy_deployment_freeze.%[5]s.id
                tenant_id = octopusdeploy_tenant.%[1]s.id
                project_id = octopusdeploy_project.%[4]s.id
                environment_id = octopusdeploy_environment.environment_0_%[5]s.id

                depends_on = [
                    octopusdeploy_tenant_project.%[1]s_project
                ]
            }`,
			tenantLocalName, tenantName, spaceLocalName, projectLocalName, localName)

		config = fmt.Sprintf("%s\n\n%s", config, tenantConfig)
	}

	return config
}

func testDeploymentFreezeExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		freezeId := s.RootModule().Resources[prefix].Primary.ID
		if _, err := deploymentfreezes.GetById(octoClient, freezeId); err != nil {
			return err
		}
		return nil
	}
}

func testDeploymentFreezeTenantExists(prefix string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		featureToggle := os.Getenv("OCTOPUS__FeatureToggles__DeploymentFreezeByTenantFeatureToggle")
		t.Logf("DeploymentFreezeByTenantFeatureToggle value: '%s'", featureToggle)

		if !ok {
			return fmt.Errorf("Resource not found: %s", prefix)
		}

		bits := strings.Split(rs.Primary.ID, ":")
		if len(bits) != 4 {
			return fmt.Errorf("Invalid ID format for deployment freeze tenant: %s", rs.Primary.ID)
		}

		freezeId := bits[0]
		tenantId := bits[1]
		projectId := bits[2]
		environmentId := bits[3]

		t.Logf("Starting tenant scope check for deployment freeze %s", freezeId)
		t.Logf("Looking for tenant: %s, project: %s, environment: %s", tenantId, projectId, environmentId)

		retryErr := resource.RetryContext(context.Background(), 2*time.Minute, func() *resource.RetryError {
			freeze, err := deploymentfreezes.GetById(octoClient, freezeId)
			if err != nil {
				t.Logf("Failed to get deployment freeze: %v", err)
				return resource.NonRetryableError(fmt.Errorf("Error getting deployment freeze: %v", err))
			}

			if freezeJSON, err := json.MarshalIndent(freeze, "", "  "); err == nil {
				t.Logf("Deployment freeze as JSON:\n%s", string(freezeJSON))
			} else {
				t.Logf("Failed to marshal freeze object to JSON: %v", err)
			}

			t.Logf("Retrieved deployment freeze with %d tenant scopes", len(freeze.TenantProjectEnvironmentScope))

			for i, scope := range freeze.TenantProjectEnvironmentScope {
				t.Logf("Scope %d - Tenant: %s, Project: %s, Environment: %s",
					i+1, scope.TenantId, scope.ProjectId, scope.EnvironmentId)
			}

			for _, scope := range freeze.TenantProjectEnvironmentScope {
				if scope.TenantId == tenantId && scope.ProjectId == projectId && scope.EnvironmentId == environmentId {
					t.Log("Found matching tenant scope in deployment freeze")
					return nil
				}
			}

			t.Log("Tenant scope not yet found, will retry...")
			return resource.RetryableError(fmt.Errorf("Tenant scope not yet found in deployment freeze (freezeId: %s)", freezeId))
		})

		if retryErr != nil {
			freeze, err := deploymentfreezes.GetById(octoClient, freezeId)
			if err != nil {
				t.Logf("Final attempt to get deployment freeze failed: %v", err)
			} else {
				t.Logf("Final state - Deployment freeze has %d tenant scopes", len(freeze.TenantProjectEnvironmentScope))
				for i, scope := range freeze.TenantProjectEnvironmentScope {
					t.Logf("Final Scope %d - Tenant: %s, Project: %s, Environment: %s",
						i+1, scope.TenantId, scope.ProjectId, scope.EnvironmentId)
				}
			}

			return fmt.Errorf("Failed to find tenant scope after retries. Error: %v", retryErr)
		}

		return nil
	}
}

func testDeploymentFreezeCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_deployment_freeze" {
			continue
		}

		feed, err := deploymentfreezes.GetById(octoClient, rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("Deployment Freeze (%s) still exists", rs.Primary.ID)
		}
	}

	os.Setenv("TF_LOG", "")
	return nil
}
