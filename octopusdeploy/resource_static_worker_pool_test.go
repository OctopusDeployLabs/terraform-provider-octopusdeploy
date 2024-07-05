package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployStaticWorkerPoolBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_static_worker_pool." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testStaticWorkerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testStaticWorkerPoolBasic(localName, name, description, isDefault, sortOrder),
				Check: resource.ComposeTestCheckFunc(
					testStaticWorkerPoolExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "is_default", strconv.FormatBool(isDefault)),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
				),
			},
		},
	})
}

func testStaticWorkerPoolBasic(
	localName string,
	name string,
	description string,
	isDefault bool,
	sortOrder int,
) string {
	return fmt.Sprintf(`resource "octopusdeploy_static_worker_pool" "%s" {
		name             = "%s"
		description      = "%s"
		is_default       = %v
		sort_order       = %v
	}`, localName, name, description, isDefault, sortOrder)
}

func testStaticWorkerPoolExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		workerPoolID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.WorkerPools.GetByID(workerPoolID); err != nil {
			return err
		}

		return nil
	}
}

func testStaticWorkerPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		workerPoolID := rs.Primary.ID
		workerPool, err := client.WorkerPools.GetByID(workerPoolID)
		if err == nil {
			if workerPool != nil {
				return fmt.Errorf("static worker pool (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

// TestWorkerPoolResource verifies that a static worker pool can be reimported with the correct settings
func TestWorkerPoolResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "15-workerpool", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "15a-workerpoolds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := workerpools.WorkerPoolsQuery{
			PartialName: "Docker",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.WorkerPools.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a worker pool called \"Docker\"")
		}
		resource := resources.Items[0].(*workerpools.StaticWorkerPool)

		if resource.WorkerPoolType != "StaticWorkerPool" {
			t.Fatal("The worker pool must be have a type of \"StaticWorkerPool\" (was \"" + resource.WorkerPoolType + "\"")
		}

		if resource.Description != "A test worker pool" {
			t.Fatal("The worker pool must be have a description of \"A test worker pool\" (was \"" + resource.Description + "\"")
		}

		if resource.SortOrder != 3 {
			t.Fatal("The worker pool must be have a sort order of \"3\" (was \"" + fmt.Sprint(resource.SortOrder) + "\"")
		}

		if resource.IsDefault {
			t.Fatal("The worker pool must be must not be the default")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "15a-workerpoolds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
