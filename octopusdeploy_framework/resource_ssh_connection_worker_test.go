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

type sshConnectionWorkerTestData struct {
	name        string
	spaceID     string
	isDisabled  bool
	host        string
	port        int64
	fingerprint string
	platform    string
}

type sshConnectionTestDependenciesData struct {
	policy  string
	pool1   string
	pool2   string
	proxy   string
	account string
}

func TestAccOctopusDeploySSHConnectionWorker(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_connection_worker." + localName
	createData := sshConnectionWorkerTestData{
		name:        acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		host:        "listening.host",
		port:        int64(acctest.RandIntRange(0, 50000)),
		fingerprint: "SHA256: " + strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
		platform:    "linux-x64",
	}
	updateData := sshConnectionWorkerTestData{
		name:        createData.name + "-updated",
		host:        "listening.host.updated",
		port:        int64(acctest.RandIntRange(0, 50000)),
		fingerprint: "SHA256: " + strconv.FormatInt(int64(acctest.RandIntRange(0, 1024)), 16),
		platform:    "osx-x64",
		isDisabled:  true,
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testSSHConnectionWorkerCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testSSHConnectionWorkerCreate(createData, localName),
				Check:  testAssertSSHConnectionWorkerCreate(createData, prefix),
			},
			{
				Config: testSSHConnectionWorkerUpdate(updateData, localName),
				Check:  testAssertSSHConnectionWorkerUpdate(updateData, prefix),
			},
		},
	})
}

func testSSHConnectionWorkerCreate(data sshConnectionWorkerTestData, localName string) string {
	source, references := testSSHConnectionWorkerDependencies(localName)

	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_ssh_connection_worker" "%s" {
			name				= "%s"
			machine_policy_id	= %s
			worker_pool_ids		= [%s]
			account_id			= %s
			host				= "%s"
			port				= %d
			fingerprint 		= "%s"
			dotnet_platform		= "%s"
		}
	`,
		source,
		localName,
		data.name,
		references.policy,
		references.pool1,
		references.account,
		data.host,
		data.port,
		data.fingerprint,
		data.platform,
	)
}

func testAssertSSHConnectionWorkerCreate(expected sshConnectionWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttrSet(prefix, "machine_policy_id"),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "1"),
		resource.TestCheckResourceAttrSet(prefix, "account_id"),
		resource.TestCheckResourceAttr(prefix, "host", expected.host),
		resource.TestCheckResourceAttr(prefix, "port", strconv.FormatInt(expected.port, 10)),
		resource.TestCheckResourceAttr(prefix, "fingerprint", expected.fingerprint),
		resource.TestCheckResourceAttr(prefix, "dotnet_platform", expected.platform),
		resource.TestCheckNoResourceAttr(prefix, "proxy_id"),
		resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
	)
}

func testSSHConnectionWorkerUpdate(data sshConnectionWorkerTestData, localName string) string {
	source, references := testSSHConnectionWorkerDependencies(localName)

	return fmt.Sprintf(`
		%s
		resource "octopusdeploy_ssh_connection_worker" "%s" {
			name				= "%s"
			machine_policy_id	= %s
			worker_pool_ids		= [%s, %s]
			account_id			= %s
			host				= "%s"
			port				= %d
			fingerprint 		= "%s"
			dotnet_platform		= "%s"
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
		references.account,
		data.host,
		data.port,
		data.fingerprint,
		data.platform,
		references.proxy,
		data.isDisabled,
	)
}

func testAssertSSHConnectionWorkerUpdate(expected sshConnectionWorkerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttrSet(prefix, "machine_policy_id"),
		resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "2"),
		resource.TestCheckResourceAttrSet(prefix, "account_id"),
		resource.TestCheckResourceAttr(prefix, "host", expected.host),
		resource.TestCheckResourceAttr(prefix, "port", strconv.FormatInt(expected.port, 10)),
		resource.TestCheckResourceAttr(prefix, "fingerprint", expected.fingerprint),
		resource.TestCheckResourceAttr(prefix, "dotnet_platform", expected.platform),
		resource.TestCheckResourceAttrSet(prefix, "proxy_id"),
		resource.TestCheckResourceAttr(prefix, "is_disabled", strconv.FormatBool(expected.isDisabled)),
	)
}

func testSSHConnectionWorkerDependencies(localName string) (string, sshConnectionTestDependenciesData) {
	policy := fmt.Sprintf("policy_%s", localName)
	pool1 := fmt.Sprintf("pool_1_%s", localName)
	pool2 := fmt.Sprintf("pool_2_%s", localName)
	proxy := fmt.Sprintf("proxy_%s", localName)
	account := fmt.Sprintf("account_%s", localName)
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

		resource "octopusdeploy_ssh_key_account" "%s" {
		  	name             = "SSH Key Pair Account"
		  	private_key_file = "[private_key_file]"
		  	username         = "[username]"
		}
		`,
		policy,
		pool1,
		pool2,
		proxy,
		account,
	)

	dependencies := sshConnectionTestDependenciesData{
		policy:  fmt.Sprintf("octopusdeploy_machine_policy.%s.id", policy),
		pool1:   fmt.Sprintf("octopusdeploy_static_worker_pool.%s.id", pool1),
		pool2:   fmt.Sprintf("octopusdeploy_static_worker_pool.%s.id", pool2),
		proxy:   fmt.Sprintf("octopusdeploy_machine_proxy.%s.id", proxy),
		account: fmt.Sprintf("octopusdeploy_ssh_key_account.%s.id", account),
	}

	return source, dependencies
}

func testSSHConnectionWorkerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_ssh_connection_worker" {
			continue
		}

		feed, err := workers.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("ssh connection worker (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
