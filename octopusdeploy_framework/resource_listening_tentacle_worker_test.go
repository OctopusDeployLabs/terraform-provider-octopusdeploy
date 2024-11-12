package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"strings"
	"testing"
)

type listeningTentacleWorkerTestData struct {
	name            string
	spaceID         string
	isDisabled      bool
	workerPoolIDs   []string
	machinePolicyID string
	uri             string
	thumbprint      string
	proxyID         string
}

func TestAccOctopusDeployListeningTentacleWorker(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_listening_tentacle_worker." + localName
	createData := listeningTentacleWorkerTestData{
		name:            acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		workerPoolIDs:   []string{acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)},
		machinePolicyID: acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
		uri:             "https://listening.test",
		thumbprint:      strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
	}
	updateData := listeningTentacleWorkerTestData{
		name:            createData.name + "-updated",
		workerPoolIDs:   append(createData.workerPoolIDs, acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)),
		machinePolicyID: acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
		uri:             "https://listening.test.updated",
		thumbprint:      strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
		isDisabled:      true,
		proxyID:         acctest.RandStringFromCharSet(8, acctest.CharSetAlpha),
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testListeningTentacleWorkerCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testListeningTentacleWorkerMandatory(createData, localName),
				Check:  testAssertListeningTentacleWorkerMandatoryAttributes(createData, prefix),
			},
			{
				Config: testListeningTentacleWorkerAll(updateData, localName),
				Check:  testAssertListeningTentacleWorkerAllAttributes(updateData, prefix),
			},
		},
	})
}

func testListeningTentacleWorkerMandatory(data listeningTentacleWorkerTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_listening_tentacle_worker" "%s" {
			name				= "%s"
			machine_policy_id	= "%s"
			worker_pool_ids		= %s
			uri					= "%s"
			thumbprint			= "%s"
		}
	`,
		localName,
		data.name,
		data.machinePolicyID,
		testSerializeWorkerPoolIdsForResource(data.workerPoolIDs),
		data.uri,
		data.thumbprint,
	)
}

func testListeningTentacleWorkerAll(data listeningTentacleWorkerTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_listening_tentacle_worker" "%s" {
			name				= "%s"
			machine_policy_id	= "%s"
			worker_pool_ids		= %s
			uri					= "%s"
			thumbprint			= "%s"
			proxy_id			= "%s"
			is_disabled			= "%v"
		}
	`,
		localName,
		data.name,
		data.machinePolicyID,
		testSerializeWorkerPoolIdsForResource(data.workerPoolIDs),
		data.uri,
		data.thumbprint,
		data.proxyID,
		data.isDisabled,
	)
}

func testAssertListeningTentacleWorkerMandatoryAttributes(expected listeningTentacleWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "machine_policy_id", expected.machinePolicyID),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids", testSerializeWorkerPoolIdsForResource(expected.workerPoolIDs)),
		resource.TestCheckResourceAttr(prefix, "uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "thumbprint", expected.thumbprint),
		resource.TestCheckNoResourceAttr(prefix, "proxy_id"),
		resource.TestCheckNoResourceAttr(prefix, "is_disabled"),
	)
}

func testAssertListeningTentacleWorkerAllAttributes(expected listeningTentacleWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "machine_policy_id", expected.machinePolicyID),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids", testSerializeWorkerPoolIdsForResource(expected.workerPoolIDs)),
		resource.TestCheckResourceAttr(prefix, "uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "thumbprint", expected.thumbprint),
		resource.TestCheckResourceAttr(prefix, "proxy_id", expected.proxyID),
		resource.TestCheckResourceAttr(prefix, "is_disabled", strconv.FormatBool(expected.isDisabled)),
	)
}

func testListeningTentacleWorkerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_listening_tentacle_worker" {
			continue
		}

		feed, err := workers.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("listening tentacle worker (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSerializeWorkerPoolIdsForResource(poolIds []string) string {
	quotedPoolIds := make([]string, len(poolIds))
	for i, poolId := range poolIds {
		quotedPoolIds[i] = fmt.Sprintf(`"%s"`, poolId)
	}

	return "[" + strings.Join(quotedPoolIds, ",") + "]"
}
