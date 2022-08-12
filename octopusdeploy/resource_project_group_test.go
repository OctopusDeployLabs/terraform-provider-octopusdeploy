package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectGroup(t *testing.T) {
	options := test.NewProjectGroupTestOptions()

	resource.Test(t, resource.TestCase{
		CheckDestroy: testProjectGroupDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

func testProjectGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		projectGroupID := rs.Primary.ID
		projectGroup, err := client.ProjectGroups.GetByID(projectGroupID)
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

		client := testAccProvider.Meta().(*client.Client)
		if _, err := client.ProjectGroups.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := client.ProjectGroups.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}
