package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployUserBasic(t *testing.T) {
	const envPrefix = "octopusdeploy_User.foo"
	const envName = "Testing one two three"
	const envDesc = "Terraform testing module User"
	const envGuided = "false"
	const envDynamic = "false"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOctopusDeployUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testUsertBasic(envName, envDesc, envGuided, envDynamic),
				Check: resource.ComposeTestCheckFunc(
					testOctopusDeployUserExists(envPrefix),
					resource.TestCheckResourceAttr(
						envPrefix, "UserName", envName),
					resource.TestCheckResourceAttr(
						envPrefix, "DisplayName", envDisplayName),
					resource.TestCheckResourceAttr(
						envPrefix, "use_guided_failure", envGuided),
					resource.TestCheckResourceAttr(
						envPrefix, "allow_dynamic_infrastructure", envDynamic),
				),
			},
		},
	})
}

func testUsertBasic(UserName, DisplayName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_User" "foo" {
			UserName           = "%s"
			displayName    = "%s"
			use_guided_failure = "%s"
			allow_dynamic_infrastructure = "%s"
		}
		`,
		UserName, DisplayName,
	)
}

func testOctopusDeployUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		return existsEnvHelper(s, client)
	}
}

func existsEnvHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.User.Get(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving User %s", err)
		}
	}
	return nil
}

func testOctopusDeployUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	return destroyEnvHelper(s, client)
}

func destroyEnvHelper(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.User.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving User %s", err)
		}
		return fmt.Errorf("User still exists")
	}
	return nil
}
