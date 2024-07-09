package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	stdslices "slices"
	"testing"
)

func TestPackageFeedCreateReleaseTriggerResources(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "52-packagefeedcreatereleasetrigger", []string{})

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatal("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	project_triggers, err := client.ProjectTriggers.GetByProjectID(resource.ID)

	if err != nil {
		t.Fatal(err.Error())
	}

	tr1Name := "My first trigger"
	tr2Name := "My second trigger"
	tr3Name := "My third trigger"

	tr1Index := stdslices.IndexFunc(project_triggers, func(t *triggers.ProjectTrigger) bool { return t.Name == tr1Name })
	tr2Index := stdslices.IndexFunc(project_triggers, func(t *triggers.ProjectTrigger) bool { return t.Name == tr2Name })
	tr3Index := stdslices.IndexFunc(project_triggers, func(t *triggers.ProjectTrigger) bool { return t.Name == tr3Name })

	if tr1Index == -1 || tr2Index == -1 || tr3Index == -1 {
		t.Fatalf("Unable to find all triggers. Expecting there to be \"%s\", \"%s\", and \"%s\".", tr1Name, tr2Name, tr3Name)
	}

	for _, triggerIndex := range []int{tr1Index, tr2Index, tr3Index} {
		if project_triggers[triggerIndex].Filter.GetFilterType() != filters.FeedFilter {
			t.Fatal("The project triggers must all be of \"FeedFilter\" type")
		}
	}

	if project_triggers[tr1Index].IsDisabled {
		t.Fatalf("The trigger \"%s\" should not be disabled", tr1Name)
	}

	if !project_triggers[tr2Index].IsDisabled {
		t.Fatalf("The trigger \"%s\" should be disabled", tr2Name)
	}

	if project_triggers[tr3Index].IsDisabled {
		t.Fatalf("The trigger \"%s\" should not be disabled", tr3Name)
	}

	tr1Filter := project_triggers[tr1Index].Filter.(*filters.FeedTriggerFilter)
	tr2Filter := project_triggers[tr2Index].Filter.(*filters.FeedTriggerFilter)
	tr3Filter := project_triggers[tr3Index].Filter.(*filters.FeedTriggerFilter)

	if len(tr1Filter.Packages) != 2 {
		t.Fatalf("The trigger \"%s\" should have 2 package references", tr1Name)
	}

	if len(tr2Filter.Packages) != 1 {
		t.Fatalf("The trigger \"%s\" should have 1 package reference", tr2Name)
	}

	if len(tr3Filter.Packages) != 3 {
		t.Fatalf("The trigger \"%s\" should have 3 package reference", tr3Name)
	}
}
