package octopusdeploy

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployDynamicWorkerPoolBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_dynamic_worker_pool." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerType := workerpools.WorkerType("UbuntuDefault")
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testDynamicWorkerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDynamicWorkerPoolBasic(localName, name, workerType, description, isDefault, sortOrder),
				Check: resource.ComposeTestCheckFunc(
					testDynamicWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "worker_type", string(workerType)),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "is_default", strconv.FormatBool(isDefault)),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
				),
			},
		},
	})
}

func testDynamicWorkerPoolBasic(
	localName string,
	name string,
	workerType workerpools.WorkerType,
	description string,
	isDefault bool,
	sortOrder int,
) string {
	return fmt.Sprintf(`resource "octopusdeploy_dynamic_worker_pool" "%s" {
		name             = "%s"
		worker_type      = "%s"
		description      = "%s"
		is_default       = %v
		sort_order       = %v
	}`, localName, name, workerType, description, isDefault, sortOrder)
}

func testDynamicWorkerPoolExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		workerPoolID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.WorkerPools.GetByID(workerPoolID); err != nil {
			return err
		}

		return nil
	}
}

func testDynamicWorkerPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		workerPoolID := rs.Primary.ID
		workerPool, err := client.WorkerPools.GetByID(workerPoolID)
		if err == nil {
			if workerPool != nil {
				return fmt.Errorf("dynamic worker pool (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
