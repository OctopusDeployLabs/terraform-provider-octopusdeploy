package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	stdslices "slices"
	"testing"
)

func TestProjectScheduledTriggerResources(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "53-scheduledprojecttrigger", []string{})

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := projects.ProjectsQuery{
			Skip: 0,
			Take: 2,
		}

		resources, err := client.Projects.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) != 2 {
			t.Fatal("There must be exactly 2 projects in the space")
		}

		nonTenantedProjectName := "Non Tenanted"
		nonTenantedProjectIndex := stdslices.IndexFunc(resources.Items, func(t *projects.Project) bool { return t.Name == nonTenantedProjectName })
		nonTenantedProject := resources.Items[nonTenantedProjectIndex]

		tenantedProjectName := "Tenanted"
		tenantedProjectIndex := stdslices.IndexFunc(resources.Items, func(t *projects.Project) bool { return t.Name == tenantedProjectName })
		tenantedProject := resources.Items[tenantedProjectIndex]

		projectTriggers, err := client.ProjectTriggers.GetAll()
		if err != nil {
			return err
		}

		nonTenantedProjectTriggersCount := 0
		tenantedProjectTriggersCount := 0

		for _, trigger := range projectTriggers {
			if trigger.ProjectID == nonTenantedProject.ID {
				nonTenantedProjectTriggersCount++
			} else if trigger.ProjectID == tenantedProject.ID {
				tenantedProjectTriggersCount++
			}
		}

		if nonTenantedProjectTriggersCount != 9 {
			t.Fatal("Non Tenanted project should have exactly 8 project triggers and 1 runbook trigger, only found: " + fmt.Sprint(nonTenantedProjectTriggersCount))
		}

		if tenantedProjectTriggersCount != 2 {
			t.Fatal("Tenanted project should have exactly 1 project trigger and 1 runbook trigger, only found: " + fmt.Sprint(tenantedProjectTriggersCount))
		}

		return nil
	})
}
