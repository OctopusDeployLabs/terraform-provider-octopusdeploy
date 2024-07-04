package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
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

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccEnvironment(environment1LocalName, environment1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment2LocalName, environment2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment3LocalName, environment3Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
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
		client := testAccProvider.Meta().(*client.Client)
		if err := existsHelperLifecycle(s, client); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
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

func existsHelperLifecycle(s *terraform.State, client *client.Client) error {
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

		client := testAccProvider.Meta().(*client.Client)
		lifecycle, err := client.Lifecycles.GetByID(rs.Primary.ID)
		if err == nil && lifecycle != nil {
			return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestLifecycleResource verifies that a lifecycle can be reimported with the correct settings
func TestLifecycleResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "17-lifecycle", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "17a-lifecycleds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := lifecycles.Query{
			PartialName: "Simple",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Lifecycles.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an environment called \"Simple\"")
		}
		resource := resources.Items[0]

		if resource.Description != "A test lifecycle" {
			t.Fatal("The lifecycle must be have a description of \"A test lifecycle\" (was \"" + resource.Description + "\")")
		}

		if resource.TentacleRetentionPolicy.QuantityToKeep != 30 {
			t.Fatal("The lifecycle must be have a tentacle retention policy of \"30\" (was \"" + fmt.Sprint(resource.TentacleRetentionPolicy.QuantityToKeep) + "\")")
		}

		if resource.TentacleRetentionPolicy.ShouldKeepForever {
			t.Fatal("The lifecycle must be have a tentacle retention not set to keep forever")
		}

		if resource.TentacleRetentionPolicy.Unit != "Items" {
			t.Fatal("The lifecycle must be have a tentacle retention unit set to \"Items\" (was \"" + resource.TentacleRetentionPolicy.Unit + "\")")
		}

		if resource.ReleaseRetentionPolicy.QuantityToKeep != 1 {
			t.Fatal("The lifecycle must be have a release retention policy of \"1\" (was \"" + fmt.Sprint(resource.ReleaseRetentionPolicy.QuantityToKeep) + "\")")
		}

		if !resource.ReleaseRetentionPolicy.ShouldKeepForever {
			t.Log("BUG: The lifecycle must be have a release retention set to keep forever (known bug - the provider creates this field as false)")
		}

		if resource.ReleaseRetentionPolicy.Unit != "Days" {
			t.Fatal("The lifecycle must be have a release retention unit set to \"Days\" (was \"" + resource.ReleaseRetentionPolicy.Unit + "\")")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "17a-lifecycleds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
