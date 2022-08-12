package octopusdeploy

import (
	"fmt"
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

	return fmt.Sprintf(testEnvironmentBasic(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
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
