package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	internaltest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccListeningTentacleDeploymentTarget(t *testing.T) {
	SkipCI(t, "Error: Missing required argument")
	options := internaltest.NewListeningTentacleDeploymentTargetTestOptions()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: internaltest.ListeningTentacleDeploymentTargetConfiguration(options),
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

// TestListeningTargetResource verifies that a listening machine can be reimported with the correct settings
func TestListeningTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "31-listeningtarget", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "31a-listeningtargetds"), newSpaceId, []string{})

	if err != nil {
		t.Log("BUG: listening targets data sources don't appear to work")
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := machines.MachinesQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Machines.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a machine called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.URI != "https://tentacle/" {
		t.Fatal("The machine must have a Uri of \"https://tentacle/\" (was \"" + resource.URI + "\")")
	}

	if resource.Thumbprint != "55E05FD1B0F76E60F6DA103988056CE695685FD1" {
		t.Fatal("The machine must have a Thumbprint of \"55E05FD1B0F76E60F6DA103988056CE695685FD1\" (was \"" + resource.Thumbprint + "\")")
	}

	if len(resource.Roles) != 1 {
		t.Fatal("The machine must have 1 role")
	}

	if resource.Roles[0] != "vm" {
		t.Fatal("The machine must have a role of \"vm\" (was \"" + resource.Roles[0] + "\")")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}
}
