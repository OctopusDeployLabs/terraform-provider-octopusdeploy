package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOctopusDeployEnvironmentBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_environment." + localName

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccEnvironmentCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentExists(prefix),
					resource.TestCheckResourceAttr(prefix, "allow_dynamic_infrastructure", strconv.FormatBool(allowDynamicInfrastructure)),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "sort_order", strconv.Itoa(sortOrder)),
					resource.TestCheckResourceAttr(prefix, "use_guided_failure", strconv.FormatBool(useGuidedFailure)),
				),
				Config: testAccEnvironment(localName, name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure),
			},
		},
	})
}

func TestAccOctopusDeployEnvironmentMinimum(t *testing.T) {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_environment." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccEnvironmentCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				Config: testAccEnvironment(localName, name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure),
			},
		},
	})
}

func testAccEnvironment(localName string, name string, description string, allowDynamicInfrastructure bool, sortOrder int, useGuidedFailure bool) string {
	return fmt.Sprintf(`resource "octopusdeploy_environment" "%s" {
		allow_dynamic_infrastructure = "%v"
		description                  = "%s"
		name                         = "%s"
		sort_order                   = %v
		use_guided_failure           = "%v"
	}`, localName, allowDynamicInfrastructure, description, name, sortOrder, useGuidedFailure)
}

func testAccEnvironmentExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		environmentID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Environments.GetByID(environmentID); err != nil {
			return err
		}

		return nil
	}
}

func testAccEnvironmentCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_environment" {
			continue
		}

		if environment, err := client.Environments.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("environment (%s) still exists", environment.GetID())
		}
	}

	return nil
}

// TestEnvironmentResource verifies that an environment can be reimported with the correct settings
func TestEnvironmentResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "16-environment", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "16a-environmentlookup"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := environments.EnvironmentsQuery{
			PartialName: "Development",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Environments.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an environment called \"Development\"")
		}
		resource := resources.Items[0]

		if resource.Description != "A test environment" {
			t.Fatal("The environment must be have a description of \"A test environment\" (was \"" + resource.Description + "\"")
		}

		if !resource.AllowDynamicInfrastructure {
			t.Fatal("The environment must have dynamic infrastructure enabled.")
		}

		if resource.UseGuidedFailure {
			t.Fatal("The environment must not have guided failure enabled.")
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "16a-environmentlookup"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
