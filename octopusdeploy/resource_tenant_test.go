package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTenantBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_tenant." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccTenantCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
				Config: testAccTenantBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", newDescription),
				),
				Config: testAccTenantBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, localName, name, newDescription),
			},
		},
	})
}

func testAccTenantBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, environmentLocalName string, environmentName string, localName string, name string, description string) string {
	allowDynamicInfrastructure := false
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, environmentDescription, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
	resource "octopusdeploy_tenant" "%s" {
		description = "%s"
		name        = "%s"

		project_environment {
		  project_id   = "${octopusdeploy_project.%s.id}"
		  environments = ["${octopusdeploy_environment.%s.id}"]
		}
	}`, localName, description, name, projectLocalName, environmentLocalName)
}

func testTenantExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("Not found: %s", prefix)
		}

		client := testAccProvider.Meta().(*client.Client)
		if _, err := client.Tenants.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTenantCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant" {
			continue
		}

		if tenant, err := client.Tenants.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("tenant (%s) still exists", tenant.GetID())
		}
	}

	return nil
}

// TestTenantsResource verifies that a git credential can be reimported with the correct settings
func TestTenantsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "24-tenants", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "24a-tenantsds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := tenants.TenantsQuery{
			PartialName: "Team A",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Tenants.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a tenant called \"Team A\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test tenant" {
			t.Fatal("The tenant must be have a description of \"tTest tenant\" (was \"" + resource.Description + "\")")
		}

		if len(resource.TenantTags) != 2 {
			t.Fatal("The tenant must have two tags")
		}

		if len(resource.ProjectEnvironments) != 1 {
			t.Fatal("The tenant must have one project environment")
		}

		for _, u := range resource.ProjectEnvironments {
			if len(u) != 3 {
				t.Fatal("The tenant must have be linked to three environments")
			}
		}

		// Verify the environment data lookups work
		tagsets, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "24a-tenantsds"), "tagsets")

		if err != nil {
			return err
		}

		if tagsets == "" {
			t.Fatal("The tagset lookup failed.")
		}

		tenants, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "24a-tenantsds"), "tenants_lookup")

		if err != nil {
			return err
		}

		if tenants != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + tenants + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}
