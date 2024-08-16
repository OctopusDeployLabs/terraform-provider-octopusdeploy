package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	stdslices "slices"
	"testing"
)

func TestKubernetesAgentWorkerResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{
		CustomEnvironment: map[string]string{
			"OCTOPUS__FeatureToggles__KubernetesAgentAsWorkerFeatureToggle": "true",
		},
	}
	_, err := testFramework.Act(t, octoContainer, "../terraform", "58-kubernetesagentworker", []string{})
	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	query := machines.WorkersQuery{
		CommunicationStyles: []string{"KubernetesTentacle"},
		Skip:                0,
		Take:                3,
	}

	resources, err := octoClient.Workers.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) != 2 {
		t.Fatalf("Space must have two workers (both KubernetesTentacles), instead found %v", resources.Items)
	}

	optionalAgentName := "minimum-agent"
	optionalAgentIndex := stdslices.IndexFunc(resources.Items, func(t *machines.Worker) bool { return t.Name == optionalAgentName })
	optionalAgentWorker := resources.Items[optionalAgentIndex]
	optionalAgentEndpoint := optionalAgentWorker.Endpoint.(*machines.KubernetesTentacleEndpoint)
	if optionalAgentWorker.IsDisabled {
		t.Fatalf("Expected  \"%s\" to be enabled", optionalAgentName)
	}

	if optionalAgentEndpoint.UpgradeLocked {
		t.Fatalf("Expected  \"%s\" to not be upgrade locked", optionalAgentName)
	}

	fullAgentName := "agent-with-optionals"
	fullAgentIndex := stdslices.IndexFunc(resources.Items, func(t *machines.Worker) bool { return t.Name == fullAgentName })
	fullAgentWorker := resources.Items[fullAgentIndex]

	if !fullAgentWorker.IsDisabled {
		t.Fatalf("Expected  \"%s\" to be disabled", fullAgentName)
	}

	fullAgentEndpoint := fullAgentWorker.Endpoint.(*machines.KubernetesTentacleEndpoint)
	if !fullAgentEndpoint.UpgradeLocked {
		t.Fatalf("Expected  \"%s\" to be upgrade locked", fullAgentName)
	}

	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "58-kubernetesagentworker"), "data_lookup")

	if err != nil {
		t.Fatal("Failed to query for created k8s workers")
	}

	if len(lookup) > 5 {
		t.Fatal("Failed to query for created k8s workers")
	}
}
