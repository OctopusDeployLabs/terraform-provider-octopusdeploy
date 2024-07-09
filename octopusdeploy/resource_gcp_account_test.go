package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGcpAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcp_account." + localName

	jsonKey := acctest.RandString(acctest.RandIntRange(20, 255))
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "json_key", jsonKey),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
				),
				Config: testGcpAccountBasic(localName, name, description, jsonKey, tenantedDeploymentParticipation),
			},
		},
	})
}

func testGcpAccountBasic(localName string, name string, description string, jsonKey string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcp_account" "%s" {
		description = "%s"
		json_key = "%s"
		name = "%s"
		tenanted_deployment_participation = "%s"
	}

	data "octopusdeploy_accounts" "test" {
		ids = [octopusdeploy_gcp_account.%s.id]
	}`, localName, description, jsonKey, name, tenantedDeploymentParticipation, localName)
}

func testGcpAccount(localName string, name string, jsonKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcp_account" "%s" {
		json_key = "%s"
		name = "%s"
	}`, localName, jsonKey, name)
}
