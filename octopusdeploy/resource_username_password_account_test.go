package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUsernamePasswordBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_username_password_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := testUsernamePasswordBasic(localName, description, name, username, password, tenantedDeploymentParticipation)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
				ResourceName: resourceName,
			},
		},
	})
}

func testUsernamePasswordBasic(localName string, description string, name string, username string, password string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_username_password_account" "%s" {
		description                       = "%s"
		name                              = "%s"
		password                          = "%s"
		tenanted_deployment_participation = "%s"
		username                          = "%s"
	}`, localName, description, name, password, tenantedDeploymentParticipation, username)
}
