package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func (suite *IntegrationTestSuite) TestAccOctopusDeployEnvironmentBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_environment." + localName

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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

func (suite *IntegrationTestSuite) TestAccOctopusDeployEnvironmentMinimum() {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_environment." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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
		environmentID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Environments.GetByID(environmentID); err != nil {
			return err
		}

		return nil
	}
}

func testAccEnvironmentCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_environment" {
			continue
		}

		if environment, err := octoClient.Environments.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("environment (%s) still exists", environment.GetID())
		}
	}

	return nil
}

// TestEnvironmentResource verifies that an environment can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestEnvironmentResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "16-environment", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "16a-environmentlookup"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := environments.EnvironmentsQuery{
		PartialName: "Development",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Environments.Get(query)
	if err != nil {
		t.Fatal(err.Error())
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
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
