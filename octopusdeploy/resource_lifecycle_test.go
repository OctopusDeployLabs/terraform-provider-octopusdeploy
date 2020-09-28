package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
						terraformNamePrefix, constName, lifecycleName),
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
						terraformNamePrefix, constName, lifecycleName),
				),
			},
			// create update it with a description
			{
				Config: testAccLifecycleWithDescription(lifecycleName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, lifecycleName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constDescription, description),
				),
			},
			// update again by remove its description
			{
				Config: testAccLifecycleBasic(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(terraformNamePrefix),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constName, lifecycleName),
					resource.TestCheckResourceAttr(
						terraformNamePrefix, constDescription, ""),
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
						terraformNamePrefix, constName, "Funky Lifecycle"),
				),
			},
		},
	})
}

func testAccLifecycleBasic(name string) string {
	return fmt.Sprintf(`
		resource constOctopusDeployLifecycle "foo" {
			name           = "%s"
		  }
		`,
		name,
	)
}
func testAccLifecycleWithDescription(name, description string) string {
	return fmt.Sprintf(`
		resource constOctopusDeployLifecycle "foo" {
			name           = "%s"
			description    = "%s"
		  }
		`,
		name, description,
	)
}

func testAccLifecycleComplex() string {
	return `
        resource constOctopusDeployEnvironment "Env1" {
           name =  "LifecycleTestEnv1"        
        }

        resource constOctopusDeployEnvironment "Env2" {
           name =  "LifecycleTestEnv2"
        }

 		resource constOctopusDeployEnvironment "Env3" {
           name =  "LifecycleTestEnv3"
        }

        resource constOctopusDeployLifecycle "foo" {
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
	client := testAccProvider.Meta().(*client.Client)

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
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelperLifecycle(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckOctopusDeployLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		resourceList, err := client.Lifecycles.GetByPartialName(name)
		if err != nil {
			return err
		}

		resource := resourceList[0]

		if len(resource.Phases) != expected {
			return fmt.Errorf("Lifecycle has %d phases instead of the expected %d", len(resource.Phases), expected)
		}

		return nil
	}
}
func destroyHelperLifecycle(s *terraform.State, apiClient *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := apiClient.Lifecycles.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving lifecycle %s", err)
		}
		return fmt.Errorf("lifecycle still exists")
	}
	return nil
}

func existsHelperLifecycle(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == constOctopusDeployLifecycle {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}
