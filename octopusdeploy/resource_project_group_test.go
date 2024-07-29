package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccProjectGroup() {
	options := test.NewProjectGroupTestOptions()
	t := suite.T()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testProjectGroupDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testProjectGroupExists(options.QualifiedName),
					resource.TestCheckResourceAttr(options.QualifiedName, "description", options.Resource.Description),
					resource.TestCheckResourceAttr(options.QualifiedName, "name", options.Resource.Name),
				),
				Config: test.ProjectGroupConfiguration(options),
			},
		},
	})
}

func testAccProjectGroup(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testProjectGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		projectGroupID := rs.Primary.ID
		projectGroup, err := octoClient.ProjectGroups.GetByID(projectGroupID)
		if err == nil {
			if projectGroup != nil {
				return fmt.Errorf("project group (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testProjectGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}
