package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strings"
	"testing"
	"time"
)

func TestNewDeploymentFreezeResource(t *testing.T) {
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
				Config: testDeploymentFreezeBasic(localName, name, start, end, spaceName, []string{environmentName1}, projectName, projectGroupName, lifecycleName, tenantName, false, false),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"1"),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", updatedEnd)),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, []string{environmentName1, environmentName2}, projectName, projectGroupName, lifecycleName, tenantName, false, false),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					testDeploymentFreezeTenantExists(fmt.Sprintf("octopusdeploy_deployment_freeze_tenant.tenant_%s", localName))),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, []string{environmentName1, environmentName2}, projectName, projectGroupName, lifecycleName, tenantName, true, true),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testDeploymentFreezeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name+"1"),
					resource.TestCheckResourceAttr(resourceName, "start", start),
					resource.TestCheckResourceAttr(resourceName, "end", updatedEnd),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.type", "Weekly"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.unit", "24"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.end_type", "AfterOccurrences"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.end_after_occurrences", "5"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.days_of_week.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.days_of_week.0", "Monday"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.days_of_week.1", "Wednesday"),
					resource.TestCheckResourceAttr(resourceName, "recurring_schedule.days_of_week.2", "Friday")),
				Config: testDeploymentFreezeBasic(localName, name+"1", start, updatedEnd, spaceName, []string{environmentName1, environmentName2}, projectName, projectGroupName, lifecycleName, tenantName, true, true),
			},
		},
	})
}

func testDeploymentFreezeBasic(localName string, freezeName string, start string, end string, spaceName string, environments []string, projectName string, projectGroupName string, lifecycleName string, tenantName string, includeTenant bool, includeRecurringSchedule bool) string {
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

	freezeConfig := fmt.Sprintf(`
        resource "octopusdeploy_deployment_freeze" "%s" {
            name = "%s"
            start = "%s"
            end = "%s"`, localName, freezeName, start, end)

	if includeRecurringSchedule {
		freezeConfig += `
            recurring_schedule = {
			  	type = "Weekly"          
                unit = 24
                end_type = "AfterOccurrences"
                end_after_occurrences = 5
                days_of_week = ["Monday", "Wednesday", "Friday"]
            }`
	}

	freezeConfig += `
        }`

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

        %s

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
		freezeConfig,
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

func testDeploymentFreezeTenantExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		bits := strings.Split(rs.Primary.ID, ":")
		if len(bits) != 4 {
			return fmt.Errorf("Invalid ID format for deployment freeze tenant: %s", rs.Primary.ID)
		}

		freezeId := bits[0]
		tenantId := bits[1]
		projectId := bits[2]
		environmentId := bits[3]

		freeze, err := deploymentfreezes.GetById(octoClient, freezeId)
		if err != nil {
			return err
		}

		for _, scope := range freeze.TenantProjectEnvironmentScope {
			if scope.TenantId == tenantId && scope.ProjectId == projectId && scope.EnvironmentId == environmentId {
				return nil
			}
		}

		return fmt.Errorf("Tenant scope not found in deployment freeze")
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
	return nil
}
