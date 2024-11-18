package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceWorkers(t *testing.T) {
	localName := acctest.RandStringFromCharSet(50, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_workers.%s", localName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: createTestAccDataSourceWorkerResources(),
			},
			{
				Config: createTestAccDataSourceWorkersEmpty(localName),
				Check:  testAssertDataSourceWorkersEmpty(prefix),
			},
			{
				Config: createTestAccDataSourceSSHWorkers(localName),
				Check:  testAssertDataSourceSSHWorkers(prefix),
			},
			{
				Config: createTestAccDataSourceListeningWorkers(localName),
				Check:  testAssertDataSourceListeningWorkers(prefix),
			},
		},
	})
}

func createTestAccDataSourceWorkersEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_workers" "%s" {}`, localName)
}

func testAssertDataSourceWorkersEmpty(prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAssertWorkersDataSourceID(prefix),
		resource.TestCheckResourceAttr(prefix, "workers.#", "2"),
	)
}

func createTestAccDataSourceSSHWorkers(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_workers" "%s" {
		communication_styles = ["Ssh"]
	}`, localName)
}

func testAssertDataSourceSSHWorkers(prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAssertWorkersDataSourceID(prefix),
		resource.TestCheckResourceAttr(prefix, "workers.#", "1"),
		resource.TestCheckResourceAttr(prefix, "workers[0].name", "First SSH worker"),
		resource.TestCheckResourceAttr(prefix, "workers[0].host", "test.domain"),
		resource.TestCheckResourceAttr(prefix, "workers[0].port", "4201"),
		resource.TestCheckResourceAttr(prefix, "workers[0].fingerprint", "SHA256: 1234abcdef56789"),
		resource.TestCheckResourceAttr(prefix, "workers[0].dotnet_platform", "linux-x64"),
	)
}

func createTestAccDataSourceListeningWorkers(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_workers" "%s" {
		communication_styles = ["TentaclePassive"]
	}`, localName)
}

func testAssertDataSourceListeningWorkers(prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAssertWorkersDataSourceID(prefix),
		resource.TestCheckResourceAttr(prefix, "workers.#", "1"),
		resource.TestCheckResourceAttr(prefix, "workers[0].name", "First SSH worker"),
		resource.TestCheckResourceAttr(prefix, "workers[0].uri", "https://domain.test/"),
		resource.TestCheckResourceAttr(prefix, "workers[0].thumbprint", "absdef"),
	)
}

func testAssertWorkersDataSourceID(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		dataSource, ok := all[prefix]
		if !ok {
			return fmt.Errorf("cannot find Workers data source: %s", prefix)
		}

		if dataSource.Primary.ID == "" {
			return fmt.Errorf("snapshot Workers source ID not set")
		}
		return nil
	}
}

func createTestAccDataSourceWorkerResources() string {
	return `resource "octopusdeploy_machine_policy" "policy_1" {
		  	name = "Machine Policy One"
		}

		resource "octopusdeploy_static_worker_pool" "pool_1" {
			name		= "Worker Pool One"
			description	= "First pool of listening workers"
			sort_order	= 99
		}

		resource "octopusdeploy_ssh_key_account" "account_1" {
		  	name             = "SSH Key Pair Account"
		  	private_key_file = "[private_key_file]"
		  	username         = "[username]"
		}

		resource "octopusdeploy_ssh_connection_worker" "worker_1" {
			name				= "First SSH worker"
			machine_policy_id	= octopusdeploy_machine_policy.policy_1.id
			worker_pool_ids		= [octopusdeploy_static_worker_pool.pool_1.id]
			account_id			= octopusdeploy_ssh_key_account.account_1.id
			host				= "test.domain"
			port				= 4201
			fingerprint 		= "SHA256: 1234abcdef56789"
			dotnet_platform		= "linux-x64"
		}

		resource "octopusdeploy_listening_tentacle_worker" "worker_2" {
			name				= "Second listening worker"
			machine_policy_id	= octopusdeploy_machine_policy.policy_1.id
			worker_pool_ids		= [octopusdeploy_static_worker_pool.pool_1.id]
			uri					= "https://domain.test/"
			thumbprint			= "abcdef"
		}
`
}
