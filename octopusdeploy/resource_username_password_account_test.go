package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUsernamePasswordBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployUsernamePasswordAccount + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUsernamePasswordBasic(localName, name, username, password, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constPassword, password),
					resource.TestCheckResourceAttr(prefix, constTenantedDeploymentParticipation, string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, constUsername, username),
				),
			},
		},
	})
}

func testUsernamePasswordBasic(localName string, name string, username string, password string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
		password = "%s"
		tenanted_deployment_participation = "%s"
		username = "%s"
	}`, constOctopusDeployUsernamePasswordAccount, localName, name, password, tenantedDeploymentParticipation, username)
}
