package octopusdeploy

import (
	"fmt"
	"testing"

	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccListeningTentacleDeploymentTarget(t *testing.T) {
	internalTest.SkipCI(t, "Error: Missing required argument")
	options := internalTest.NewListeningTentacleDeploymentTargetTestOptions()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: internalTest.ListeningTentacleDeploymentTargetConfiguration(options),
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(options.ResourceName),
				),
			},
			{
				ResourceName:      options.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// func TestAccListeningTentacleDeploymentTargetSchemaValidation(t *testing.T) {
// 	options := test.NewTestOptions()
// 	resourceName := "octopusdeploy_listening_tentacle_deployment_target." + options.LocalName

// 	resource.Test(t, resource.TestCase{
// 		CheckDestroy: testAccListeningTentacleDeploymentTargetCheckDestroy,
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccListeningTentacleDeploymentTargetSchemaValidation(options),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccListeningTentacleDeploymentTargetExists(resourceName),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccListeningTentacleDeploymentTargetBasic(options *test.TestOptions) string {
// 	allowDynamicInfrastructure := false
// 	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	sortOrder := acctest.RandIntRange(0, 10)
// 	thumbprint := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	useGuidedFailure := false

// 	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "default" {
// 		partial_name = "Default Machine Policy"
// 	}`+"\n"+
// 		testEnvironmentBasic(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
// 	resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
// 		environments                      = [octopusdeploy_environment.%s.id]
// 		is_disabled                       = true
// 		machine_policy_id                 = data.octopusdeploy_machine_policies.default.machine_policies[0].id
// 		name                              = "%s"
// 		roles                             = ["Prod"]
// 		tentacle_url                      = "https://example.com:1234/"
// 		thumbprint                        = "%s"
// 	  }`, options.LocalName, environmentLocalName, options.Name, thumbprint)
// }

// func testAccListeningTentacleDeploymentTargetSchemaValidation(options *test.TestOptions) string {
// 	allowDynamicInfrastructure := false
// 	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	sortOrder := acctest.RandIntRange(0, 10)
// 	thumbprint := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
// 	useGuidedFailure := false

// 	return fmt.Sprintf(testEnvironmentBasic(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
// 		`resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
// 			environments = [octopusdeploy_environment.%s.id]
// 			name         = "%s"
// 			roles        = ["something"]
// 			tentacle_url = "https://example.com/"
// 			thumbprint   = "%s"
// 	  	}`, options.LocalName, environmentLocalName, options.Name, thumbprint)
// }

func testAccListeningTentacleDeploymentTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[resourceName].Primary.ID
		if _, err := octoClient.Machines.GetByID(deploymentTargetID); err != nil {
			return fmt.Errorf("error retrieving deployment target: %s", err)
		}

		return nil
	}
}

func testAccListeningTentacleDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_listening_tentacle_deployment_target" {
			continue
		}

		_, err := octoClient.Machines.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("deployment target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
