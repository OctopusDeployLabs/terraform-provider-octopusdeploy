package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"path/filepath"
	"strconv"
)

func (suite *IntegrationTestSuite) TestAccOctopusDeployStaticWorkerPoolBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_static_worker_pool." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	isDefault := false
	sortOrder := acctest.RandIntRange(50, 100)
	t := suite.T()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testStaticWorkerPoolDestroy,
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
		workerPoolID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.WorkerPools.GetByID(workerPoolID); err != nil {
			return err
		}

		return nil
	}
}

func testStaticWorkerPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		workerPoolID := rs.Primary.ID
		workerPool, err := octoClient.WorkerPools.GetByID(workerPoolID)
		if err == nil {
			if workerPool != nil {
				return fmt.Errorf("static worker pool (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

// TestWorkerPoolResource verifies that a static worker pool can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestWorkerPoolResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "15-workerpool", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "15a-workerpoolds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := workerpools.WorkerPoolsQuery{
		PartialName: "Docker",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.WorkerPools.Get(query)
	if err != nil {
		t.Fatal(err.Error())
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
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
