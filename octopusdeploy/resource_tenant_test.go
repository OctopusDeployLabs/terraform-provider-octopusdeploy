package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"path/filepath"
)

func (suite *IntegrationTestSuite) TestAccTenantBasic() {
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
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
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

		if _, err := octoClient.Tenants.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTenantCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant" {
			continue
		}

		if tenant, err := octoClient.Tenants.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("tenant (%s) still exists", tenant.GetID())
		}
	}

	return nil
}

// TestTenantsResource verifies that a git credential can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestTenantsResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "24-tenants", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "24a-tenantsds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := tenants.TenantsQuery{
		PartialName: "Team A",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Tenants.Get(query)
	if err != nil {
		t.Fatal(err.Error())
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
		t.Fatal(err.Error())
	}

	if tagsets == "" {
		t.Fatal("The tagset lookup failed.")
	}

	tenants, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "24a-tenantsds"), "tenants_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if tenants != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + tenants + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
