package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployLifecycleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployLifecycle + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
				Config: testAccLifecycleBasic(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployLifecycleWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployLifecycle + "." + localName

	description := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testLifecycleDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// create lifecycle with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
				Config: testAccLifecycleBasic(localName, name),
			},
			// update lifecycle with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constDescription, description),
				),
				Config: testAccLifecycleWithDescription(localName, name, description),
			},
			// update lifecycle by removing its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constDescription, ""),
				),
				Config: testAccLifecycleBasic(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployLifecycleComplex(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployLifecycle + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccCheckOctopusDeployLifecycleDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLifecycleExists(prefix),
					testAccCheckOctopusDeployLifecyclePhaseCount(name, 2),
					resource.TestCheckResourceAttr(prefix, constName, name),
				),
				Config: testAccLifecycleComplex(localName, name),
			},
		},
	})
}

func testAccLifecycleBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
	}`, constOctopusDeployLifecycle, localName, name)
}

func testAccLifecycleWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		description = "%s"
		name        = "%s"
	}`, constOctopusDeployLifecycle, localName, description, name)
}

func testAccLifecycleComplex(localName string, name string) string {
	environment1LocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	environment1Name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	environment2LocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	environment2Name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	environment3LocalName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	environment3Name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	return fmt.Sprintf(testEnvironmentMinimum(environment1LocalName, environment1Name)+"\n"+
		testEnvironmentMinimum(environment2LocalName, environment2Name)+"\n"+
		testEnvironmentMinimum(environment3LocalName, environment3Name)+"\n"+
		`resource "%s" "%s" {
			name        = "%s"
			description = "Funky Lifecycle description"
			release_retention_policy {
				unit             = "Days"
				quantity_to_keep = 2
			}
			tentacle_retention_policy {
				unit             = "Days"
				quantity_to_keep = 1
			}
			phase {
				automatic_deployment_targets          = ["${octopusdeploy_environment.%s.id}"]
				is_optional_phase                     = true
				minimum_environments_before_promotion = 2
				name                                  = "P1"
				optional_deployment_targets           = ["${octopusdeploy_environment.%s.id}"]
			}
			phase {
				name = "P2"
			}
	}`, constOctopusDeployLifecycle, localName, name, environment2LocalName, environment3LocalName)
}

func testAccCheckOctopusDeployLifecycleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)

	if err := destroyHelperLifecycle(s, client); err != nil {
		return err
	}
	if err := testEnvironmentDestroy(s); err != nil {
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
func destroyHelperLifecycle(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
			return fmt.Errorf("Received an error retrieving lifecycle %s", err)
		}
		return fmt.Errorf("lifecycle still exists")
	}
	return nil
}

func existsHelperLifecycle(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == constOctopusDeployLifecycle {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("received an error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}

func testLifecycleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		lifecycleID := rs.Primary.ID
		lifecycle, err := client.Lifecycles.GetByID(lifecycleID)
		if err == nil {
			if lifecycle != nil {
				return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
