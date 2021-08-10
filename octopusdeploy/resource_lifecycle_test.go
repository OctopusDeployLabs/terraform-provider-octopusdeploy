package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLifecycleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccLifecycleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
		},
	})
}

func TestAccLifecycleWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccLifecycleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// create lifecycle with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
			// update lifecycle with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycleWithDescription(localName, name, description),
			},
			// update lifecycle by removing its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
		},
	})
}

func TestAccLifecycleComplex(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccLifecycleCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "2"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
					testAccCheckLifecyclePhaseCount(name, 2),
				),
				Config: testAccLifecycleComplex(localName, name),
			},
		},
	})
}

func testAccLifecycle(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccLifecycleWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		description = "%s"
		name        = "%s"
	}`, localName, description, name)
}

func testAccLifecycleComplex(localName string, name string) string {
	environment1LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment1Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testAccEnvironment(environment1LocalName, environment1Name)+"\n"+
		testAccEnvironment(environment2LocalName, environment2Name)+"\n"+
		testAccEnvironment(environment3LocalName, environment3Name)+"\n"+
		`resource "octopusdeploy_lifecycle" "%s" {
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
	}`, localName, name, environment2LocalName, environment3LocalName)
}

func testAccCheckLifecycleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		if err := existsHelperLifecycle(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		resourceList, err := client.Lifecycles.GetByPartialName(name)
		if err != nil {
			return err
		}

		resource := resourceList[0]

		if len(resource.Phases) != expected {
			return fmt.Errorf("lifecycle has %d phases instead of the expected %d", len(resource.Phases), expected)
		}

		return nil
	}
}

func existsHelperLifecycle(s *terraform.State, client *octopusdeploy.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_lifecycle" {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}

func testAccLifecycleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_lifecycle" {
			continue
		}

		client := testAccProvider.Meta().(*octopusdeploy.Client)
		lifecycle, err := client.Lifecycles.GetByID(rs.Primary.ID)
		if err == nil && lifecycle != nil {
			return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
