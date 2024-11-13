package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

type listeningTentacleWorkerTestData struct {
	name       string
	spaceID    string
	isDisabled bool
	uri        string
	thumbprint string
}

type listeningTentacleWorkerTestDependenciesData struct {
	policy string
	pool1  string
	pool2  string
	proxy  string
}

func TestAccOctopusDeployListeningTentacleWorker(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_listening_tentacle_worker." + localName
	createData := listeningTentacleWorkerTestData{
		name:       acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		uri:        "https://listening.test/",
		thumbprint: strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
	}
	updateData := listeningTentacleWorkerTestData{
		name:       createData.name + "-updated",
		uri:        "https://listening.test.updated/",
		thumbprint: strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
		isDisabled: true,
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testListeningTentacleWorkerCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testListeningTentacleWorkerCreate(createData, localName),
				Check:  testAssertListeningTentacleWorkerCreate(createData, prefix),
			},
			{
				Config: testListeningTentacleWorkerUpdate(updateData, localName),
				Check:  testAssertListeningTentacleWorkerUpdate(updateData, prefix),
			},
		},
	})
}

func testListeningTentacleWorkerCreate(data listeningTentacleWorkerTestData, localName string) string {
	source, references := testListeningTentacleWorkerDependencies(localName)

	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_listening_tentacle_worker" "%s" {
			name				= "%s"
			machine_policy_id	= %s
			worker_pool_ids		= [%s]
			uri					= "%s"
			thumbprint			= "%s"
		}
	`,
		source,
		localName,
		data.name,
		references.policy,
		references.pool1,
		data.uri,
		data.thumbprint,
	)
}

func testAssertListeningTentacleWorkerCreate(expected listeningTentacleWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttrSet(prefix, "machine_policy_id"),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "1"),
		resource.TestCheckResourceAttr(prefix, "uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "thumbprint", expected.thumbprint),
		resource.TestCheckNoResourceAttr(prefix, "proxy_id"),
		resource.TestCheckNoResourceAttr(prefix, "is_disabled"),
	)
}

func testListeningTentacleWorkerUpdate(data listeningTentacleWorkerTestData, localName string) string {
	source, references := testListeningTentacleWorkerDependencies(localName)

	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_listening_tentacle_worker" "%s" {
			name				= "%s"
			machine_policy_id	= %s
			worker_pool_ids		= [%s, %s]
			uri					= "%s"
			thumbprint			= "%s"
			proxy_id			= %s
			is_disabled			= %v
		}
	`,
		source,
		localName,
		data.name,
		references.policy,
		references.pool1,
		references.pool2,
		data.uri,
		data.thumbprint,
		references.proxy,
		data.isDisabled,
	)
}

func testAssertListeningTentacleWorkerUpdate(expected listeningTentacleWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttrSet(prefix, "machine_policy_id"),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "2"),
		resource.TestCheckResourceAttr(prefix, "uri", expected.uri),
		resource.TestCheckResourceAttr(prefix, "thumbprint", expected.thumbprint),
		resource.TestCheckResourceAttrSet(prefix, "proxy_id"),
		resource.TestCheckResourceAttr(prefix, "is_disabled", strconv.FormatBool(expected.isDisabled)),
	)
}

func testListeningTentacleWorkerDependencies(localName string) (string, listeningTentacleWorkerTestDependenciesData) {
	policy := fmt.Sprintf("policy_%s", localName)
	pool1 := fmt.Sprintf("pool_1_%s", localName)
	pool2 := fmt.Sprintf("pool_2_%s", localName)
	proxy := fmt.Sprintf("proxy_%s", localName)
	source := fmt.Sprintf(`
		resource "octopusdeploy_machine_policy" "%s" {
		  	name = "Listening policy"
		}

		resource "octopusdeploy_static_worker_pool" "%s" {
			name		= "Listening poll 1"
			description	= "First pool of listening workers"
			sort_order	= 42
		}

		resource "octopusdeploy_static_worker_pool" "%s" {
			name		= "Listening poll 2"
			description	= "Second pool of listening workers"
			sort_order	= 43
		}

		resource "octopusdeploy_machine_proxy" "%s" {
			name 		= "Listening proxy"
			host 		= "localhost"
			port 		= 20034
			username	= "user_proxy"
			password 	= "secret_proxy"
		}
		`,
		policy,
		pool1,
		pool2,
		proxy,
	)

	dependencies := listeningTentacleWorkerTestDependenciesData{
		policy: fmt.Sprintf("octopusdeploy_machine_policy.%s.id", policy),
		pool1:  fmt.Sprintf("octopusdeploy_static_worker_pool.%s.id", pool1),
		pool2:  fmt.Sprintf("octopusdeploy_static_worker_pool.%s.id", pool2),
		proxy:  fmt.Sprintf("octopusdeploy_machine_proxy.%s.id", proxy),
	}

	return source, dependencies
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
