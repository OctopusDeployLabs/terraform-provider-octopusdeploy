package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTenantBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_tenant." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccUserCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
				Config: testAccTenantBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", newDescription),
				),
				Config: testAccTenantBasic(localName, name, newDescription),
			},
		},
	})
}

func testAccTenantBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant" "%s" {
		description = "%s"
		name        = "%s"

		project_environment {
		  project_id   = "Projects-5521"
		  environments = ["Environments-10124"]
		}

		project_environment {
			project_id   = "Projects-5553"
			environments = []
		  }
	  }`, localName, description, name)
}

func testTenantExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*octopusdeploy.Client)
		userID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := client.Tenants.GetByID(userID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTenantCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*octopusdeploy.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant" {
			continue
		}

		_, err := client.Tenants.GetByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("tenant (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
