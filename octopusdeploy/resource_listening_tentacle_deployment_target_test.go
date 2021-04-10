package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccListeningTentacleDeploymentTargetImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccListeningTentacleDeploymentTargetBasic(localName, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccListeningTentacleDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccListeningTentacleDeploymentTargetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(resourceName),
				),
			},
		},
	})
}

func testAccListeningTentacleDeploymentTargetBasic(localName string, name string) string {
	allowDynamicInfrastructure := false
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	thumbprint := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	useGuidedFailure := false

	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "default" {
		partial_name = "Default Machine Policy"
	}`+"\n"+
		testEnvironmentBasic(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
	resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
		environments                      = ["${octopusdeploy_environment.%s.id}"]
		is_disabled                       = true
		machine_policy_id                 = "${data.octopusdeploy_machine_policies.default.machine_policies[0].id}"
		name                              = "%s"
		roles                             = ["Prod"]
		tenanted_deployment_participation = "Untenanted"
		tentacle_url                      = "https://example.com:1234/"
		thumbprint                        = "%s"
	  }`, localName, environmentLocalName, name, thumbprint)
}

func testAccListeningTentacleDeploymentTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		deploymentTargetID := s.RootModule().Resources[resourceName].Primary.ID
		if _, err := client.Machines.GetByID(deploymentTargetID); err != nil {
			return fmt.Errorf("error retrieving deployment target: %s", err)
		}

		return nil
	}
}

func testAccListeningTentacleDeploymentTargetCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_listening_tentacle_deployment_target" {
			continue
		}

		_, err := client.Machines.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("deployment target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
