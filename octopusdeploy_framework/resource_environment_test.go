package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"strconv"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccOctopusDeployEnvironmentMinimum(t *testing.T) {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_environment." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

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

func TestAccOctopusDeployEnvironmentReplacement(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceResourceName := "octopusdeploy_space." + localName
	environmentResourceName := "octopusdeploy_environment." + localName
	clientSpaceId := octoClient.GetSpaceID()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccEnvironmentAndSpaceCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentWithSpace(localName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentExistsInSpace(environmentResourceName, spaceResourceName),
				),
			},
			{
				Config: testAccEnvironmentWithSpace(localName, clientSpaceId),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentExistsInSpace(environmentResourceName, ""),
					resource.TestCheckResourceAttr(environmentResourceName, "space_id", clientSpaceId),
				),
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
		if _, err := environments.GetByID(octoClient, octoClient.GetSpaceID(), environmentID); err != nil {
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

		if environment, err := environments.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("environment (%s) still exists", environment.GetID())
		}
	}

	return nil
}

func testAccEnvironmentWithSpace(localName string, spaceId string) string {
	environmentSpaceId := "octopusdeploy_space." + localName + ".id"
	if spaceId != "" {
		environmentSpaceId = fmt.Sprintf(`"%s"`, spaceId)
	}

	return fmt.Sprintf(`
		resource "octopusdeploy_space" "%[1]s" {
			name                        = "Replacement Space"
			description                 = "A space for environment replacement."
			space_managers_teams        = ["teams-everyone"]
		}

		resource "octopusdeploy_environment" "%[1]s" {
			name		= "Replacement"
			description	= "Replacement environment"
			space_id	= %s
		}
		`,
		localName,
		environmentSpaceId,
	)
}

func testAccEnvironmentExistsInSpace(environmentResource string, spaceResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		spaceId := octoClient.GetSpaceID()
		if spaceResource != "" {
			spaceId = s.RootModule().Resources[spaceResource].Primary.ID
		}

		environmentID := s.RootModule().Resources[environmentResource].Primary.ID
		if _, err := environments.GetByID(octoClient, spaceId, environmentID); err != nil {
			return err
		}

		return nil
	}
}

func testAccEnvironmentAndSpaceCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "octopusdeploy_environment" {
			if environment, err := environments.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
				return fmt.Errorf("environment (%s) still exists", environment.GetID())
			}
		}

		if rs.Type == "octopusdeploy_space" {
			if space, err := spaces.GetByID(octoClient, rs.Primary.ID); err == nil {
				return fmt.Errorf("space (%s) still exists", space.GetID())
			}
		}
	}

	return nil
}
