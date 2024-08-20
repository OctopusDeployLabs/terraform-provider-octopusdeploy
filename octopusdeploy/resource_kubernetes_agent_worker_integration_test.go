package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"log"
	"os"
	"path/filepath"
	stdslices "slices"
	"testing"
	"time"
)

func TestKubernetesAgentWorkerResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{
		CustomEnvironment: map[string]string{
			"OCTOPUS__FeatureToggles__KubernetesAgentAsWorkerFeatureToggle": "true",
		},
	}

	// Use separate Octopus container as this test requires a custom environment variable to be set
	octoContainer, octoClient, sqlServerContainer, network, err = testFramework.ArrangeContainer()
	if err != nil {
		log.Printf("Failed to arrange containers: (%s)", err.Error())
	}
	os.Setenv("TF_ACC", "1")

	inputVars := []string{"-var=octopus_server_58-kubernetesagentworker=" + octoContainer.URI, "-var=octopus_apikey_58-kubernetesagentworker=" + test.ApiKey}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "58-kubernetesagentworker", inputVars)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "58a-kubernetesagentworkerds"), newSpaceId, inputVars)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, key)
	query := machines.WorkersQuery{
		CommunicationStyles: []string{"KubernetesTentacle"},
		Skip:                0,
		Take:                3,
	}

	resources, err := client.Workers.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) != 2 {
		t.Fatalf("Space must have two workers (both KubernetesTentacles), instead found %v", resources.Items)
	}

	minimalAgentName := "minimum-agent"
	minimalAgentIndex := stdslices.IndexFunc(resources.Items, func(t *machines.Worker) bool { return t.Name == minimalAgentName })
	minimalAgentWorker := resources.Items[minimalAgentIndex]
	minimalAgentEndpoint := minimalAgentWorker.Endpoint.(*machines.KubernetesTentacleEndpoint)
	if minimalAgentWorker.IsDisabled {
		t.Fatalf("Expected  \"%s\" to be enabled", minimalAgentName)
	}

	if minimalAgentEndpoint.UpgradeLocked {
		t.Fatalf("Expected  \"%s\" to not be upgrade locked", minimalAgentName)
	}

	if len(minimalAgentWorker.WorkerPoolIDs) != 1 {
		t.Fatalf("Expected  \"%s\" to have one worker pool id", minimalAgentName)
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

	if len(fullAgentWorker.WorkerPoolIDs) != 2 {
		t.Fatalf("Expected  \"%s\" to have two worker pool ids", fullAgentName)
	}

	_, err = testFramework.GetOutputVariable(t, filepath.Join("../terraform", "58a-kubernetesagentworkerds"), "data_lookup_kubernetes_worker_1_id")
	_, err = testFramework.GetOutputVariable(t, filepath.Join("../terraform", "58a-kubernetesagentworkerds"), "data_lookup_kubernetes_worker_2_id")
	if err != nil {
		t.Fatal("Failed to query for created k8s workers")
	}

	ctx := context.Background()

	// Waiting for the container logs to clear.
	time.Sleep(5000 * time.Millisecond)
	err = testFramework.CleanUp(ctx, octoContainer, sqlServerContainer, network)

	if err != nil {
		log.Printf("Failed to clean up containers: (%s)", err.Error())
	}
}
