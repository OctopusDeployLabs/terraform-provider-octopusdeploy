package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDeploymentTargetImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentTargetBasic(localName, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentTargetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDeploymentTargetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "environments.0"),
					resource.TestCheckResourceAttr(resourceName, "endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "has_latest_calamari"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "is_in_process"),
					resource.TestCheckResourceAttrSet(resourceName, "machine_policy_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", "Untenanted"),
				),
			},
		},
	})
}

func testAccDeploymentTargetBasic(localName string, name string) string {
	allowDynamicInfrastructure := false
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	useGuidedFailure := false

	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "default" {
		partial_name = "Default Machine Policy"
	}`+"\n"+
		testEnvironmentBasic(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, useGuidedFailure)+"\n"+`
	resource "octopusdeploy_deployment_target" "%s" {
		environments                      = ["${octopusdeploy_environment.%s.id}"]
		is_disabled                       = true
		machine_policy_id                 = "${data.octopusdeploy_machine_policies.default.machine_policies[0].id}"
		name                              = "%s"
		roles                             = ["Prod"]
		tenanted_deployment_participation = "Untenanted"

		endpoint {
		  communication_style = "None"
		  thumbprint          = ""
		  uri                 = ""
		}
	  }`, localName, environmentLocalName, name)
}

func testAccDeploymentTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		deploymentTargetID := s.RootModule().Resources[resourceName].Primary.ID
		if _, err := client.Machines.GetByID(deploymentTargetID); err != nil {
			return fmt.Errorf("error retrieving deployment target: %s", err)
		}

		return nil
	}
}

func testAccDeploymentTargetCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_deployment_target" {
			continue
		}

		_, err := client.Machines.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("deployment target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
