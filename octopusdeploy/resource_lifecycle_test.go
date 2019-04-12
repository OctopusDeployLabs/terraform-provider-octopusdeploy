package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOctopusDeployLifecycleBasic(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_lifecycle.foo"
	const lifecycleName = "Funky Cycle"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleBasic(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", lifecycleName),
				),
			},
		},
	})
}

func TestAccOctopusDeployLifecycleWithUpdate(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_lifecycle.foo"
	const lifecycleName = "Funky Cycle"
	const description = "I am a new lifecycle description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployLifecycleDestroy,
		Steps: []resource.TestStep{
			// create projectgroup with no description
			{
				Config: testAccLifecycleBasic(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", lifecycleName),
				),
			},
			// create update it with a description
			{
				Config: testAccLifecycleWithDescription(lifecycleName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", lifecycleName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", description),
				),
			},
			// update again by remove its description
			{
				Config: testAccLifecycleBasic(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", lifecycleName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "description", ""),
				),
			},
		},
	})
}

func TestAccOctopusDeployLifecycleComplex(t *testing.T) {
	const terraformNamePrefix = "octopusdeploy_lifecycle.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOctopusDeployLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleComplex(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					testAccCheckOctopusDeployLifecyclePhaseCount("Funky Lifecycle", 2),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, "name", "Funky Lifecycle"),
				),
			},
		},
	})
}

func testAccLifecycleBasic(name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "foo" {
			name           = "%s"
		  }
		`,
		name,
	)
}
func testAccLifecycleWithDescription(name, description string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "foo" {
			name           = "%s"
			description    = "%s"
		  }
		`,
		name, description,
	)
}

func testAccLifecycleComplex() string {
	return `
        resource "octopusdeploy_environment" "Env1" {
           name =  "LifecycleTestEnv1"        
        }

        resource "octopusdeploy_environment" "Env2" {
           name =  "LifecycleTestEnv2"
        }

 		resource "octopusdeploy_environment" "Env3" {
           name =  "LifecycleTestEnv3"
        }

        resource "octopusdeploy_lifecycle" "foo" {
           name        = "Funky Lifecycle"
           description = "Funky Lifecycle description"

           release_retention_policy {
               unit            = "Items"
               quantity_to_keep = 2
           }

           tentacle_retention_policy {
               unit            = "Days"
               quantity_to_keep = 1
           }

           phase {
               name = "P1"
               minimum_environments_before_promotion = 2
               is_optional_phase = true
               automatic_deployment_targets = ["${octopusdeploy_environment.Env1.id}"]
               optional_deployment_targets = ["${octopusdeploy_environment.Env2.id}"]
           }

           phase {
               name = "P2"
           }
        }
		`
}

func testAccCheckOctopusDeployLifecycleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyHelperLifecycle(s, client); err != nil {
		return err
	}
	if err := destroyEnvHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployLifecycleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelperLifecycle(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckOctopusDeployLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		lifecycle, err := client.Lifecycle.GetByName(name)

		if err != nil {
			return err
		}

		if len(lifecycle.Phases) != expected {
			return fmt.Errorf("Lifecycle has %d phases instead of the expected %d", len(lifecycle.Phases), expected)
		}

		return nil
	}
}
func destroyHelperLifecycle(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Lifecycle.Get(r.Primary.ID); err != nil {
			if err == octopusdeploy.ErrItemNotFound {
				continue
			}
			return fmt.Errorf("Received an error retrieving lifecycle %s", err)
		}
		return fmt.Errorf("lifecycle still exists")
	}
	return nil
}

func existsHelperLifecycle(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_lifecycle" {
			if _, err := client.Lifecycle.Get(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}
