package octopusdeploy

import (
	"fmt"
	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOctopusDeployGenericOpenIDConnectAccountBasic(t *testing.T) {
	internalTest.SkipCI(t, "audience is not set on initial creation")
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_generic_openid_connect_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentMode := core.TenantedDeploymentModeTenantedOrUntenanted

	executionKeys := []string{"space"}
	audience := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccountCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.0", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "audience", audience),
				),
				Config: testGenericOpenIDConnectAccountBasic(localName, name, description, tenantedDeploymentMode, executionKeys, audience),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentMode)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.0", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "audience", audience),
				),
				Config: testGenericOpenIDConnectAccountBasic(localName, name, newDescription, tenantedDeploymentMode, executionKeys, audience),
			},
		},
	})
}

func testGenericOpenIDConnectAccountBasic(localName string, name string, description string, tenantedDeploymentParticipation core.TenantedDeploymentMode, execution_subject_keys []string, audience string) string {
	return fmt.Sprintf(`resource "octopusdeploy_generic_openid_connect_account" "%s" {
		description = "%s"
		name = "%s"
		tenanted_deployment_participation = "%s"
		execution_subject_keys = %s
		health_subject_keys = %s
		account_test_subject_keys = %s
		audience = "%s"
	}
	
	data "octopusdeploy_accounts" "test" {
		ids = [octopusdeploy_generic_openid_connect_account.%s.id]
	}`, localName, description, name, tenantedDeploymentParticipation, StringArrayToTerraformArrayFormat(execution_subject_keys), audience, localName)
}
