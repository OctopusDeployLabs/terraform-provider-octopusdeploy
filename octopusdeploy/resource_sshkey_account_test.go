package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSSHKeyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeploySSHKeyAccount + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := octopusdeploy.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSSHKeyBasic(localName, name, username, passphrase, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constPassphrase, passphrase),
					resource.TestCheckResourceAttr(prefix, constTenantedDeploymentParticipation, string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, constUsername, username),
				),
			},
		},
	})
}

func testSSHKeyBasic(localName string, name string, username string, passphrase string, tenantedDeploymentParticipation octopusdeploy.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name = "%s"
		passphrase = "%s"
		tenanted_deployment_participation = "%s"
		username = "%s"
	}`, constOctopusDeploySSHKeyAccount, localName, name, passphrase, tenantedDeploymentParticipation, username)
}
