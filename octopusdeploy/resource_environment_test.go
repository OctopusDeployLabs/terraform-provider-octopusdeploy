package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOctopusDeployEnvironmentBasic(t *testing.T) {
	const envPrefix = "octopusdeploy_environment.foo"
	const envName = "Testing one two three"
	const envDesc = "Terraform testing module environment"
	const envGuided = "false"
	const envDynamic = "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testEnvironmenttBasic(envName, envDesc, envGuided, envDynamic),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployEnvironmentExists(envPrefix),
					resource.TestCheckResourceAttr(
						envPrefix, "name", envName),
					resource.TestCheckResourceAttr(
						envPrefix, "description", envDesc),
					resource.TestCheckResourceAttr(
						envPrefix, "use_guided_failure", envGuided),
					resource.TestCheckResourceAttr(
						envPrefix, "allow_dynamic_infrastructure", envDynamic),
				),
			},
		},
	})
}

func testEnvironmenttBasic(name, description, useguided string, dynamic string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "foo" {
			name           = "%s"
			description    = "%s"
			use_guided_failure = "%s"
			allow_dynamic_infrastructure = "%s"
		}
		`,
		name, description, useguided, dynamic,
	)
}

func testOctopusDeployEnvironmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		return existsEnvHelper(s, client)
	}
}

func existsEnvHelper(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Environments.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving environment %s", err)
		}
	}
	return nil
}

func testOctopusDeployEnvironmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	return destroyEnvHelper(s, client)
}

func destroyEnvHelper(s *terraform.State, apiClient *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := apiClient.Environments.Get(r.Primary.ID); err != nil {
			if err == client.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving environment %s", err)
		}
		return fmt.Errorf("Environment still exists")
	}
	return nil
}
