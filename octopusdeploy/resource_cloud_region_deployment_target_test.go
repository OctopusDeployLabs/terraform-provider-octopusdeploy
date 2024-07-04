package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudRegionDeploymentTargetImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDeploymentTargetBasic(localName, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudRegionDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDeploymentTargetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudRegionDeploymentTargetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}

func testAccCloudRegionDeploymentTargetBasic(localName string, name string) string {
	allowDynamicInfrastructure := false
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccEnvironment(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
		resource "octopusdeploy_cloud_region_deployment_target" "%s" {
		  default_worker_pool_id = "WorkerPools-41"
		  environments           = ["${octopusdeploy_environment.%s.id}"]
		  name                   = "%s"
		  roles                  = ["Prod"]
	    }`, localName, environmentLocalName, name)
}

func testAccCloudRegionDeploymentTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		deploymentTargetID := s.RootModule().Resources[resourceName].Primary.ID
		if _, err := client.Machines.GetByID(deploymentTargetID); err != nil {
			return fmt.Errorf("error retrieving deployment target: %s", err)
		}

		return nil
	}
}

func testAccCloudRegionDeploymentTargetCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_cloud_region_deployment_target" {
			continue
		}

		_, err := client.Machines.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("deployment target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestCloudRegionTargetResource verifies that a cloud region can be reimported with the correct settings
func TestCloudRegionTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "33-cloudregiontarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "33a-cloudregiontargetds"), newSpaceId, []string{})

		if err != nil {
			t.Fatal("cloud region data source does not appear to work")
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		return nil
	})
}
