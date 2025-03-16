package octopusdeploy_framework

// TODO: Add back once supported in octopusdeploy/octopusdeploy docker hub

//import (
//	"fmt"
//	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
//	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"strings"
//	"testing"
//)
//
//func TestAccGenericOidcAccountBasic(t *testing.T) {
//	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	resourceName := "octopusdeploy_generic_oidc_account." + localName
//
//	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted
//
//	executionKeys := []string{"space"}
//	audience := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	updatedAudience := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//
//	config := testGenericOidcAccountBasic(localName, name, description, tenantedDeploymentParticipation, executionKeys, audience)
//	updateConfig := testGenericOidcAccountBasic(localName, name, description, tenantedDeploymentParticipation, executionKeys, updatedAudience)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			{
//				Config: config,
//				Check: resource.ComposeTestCheckFunc(
//					testAccountExists(resourceName),
//					resource.TestCheckResourceAttr(resourceName, "description", description),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "name", name),
//					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
//					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
//					resource.TestCheckResourceAttr(resourceName, "execution_subject_keys.0", executionKeys[0]),
//					resource.TestCheckResourceAttr(resourceName, "audience", audience),
//				),
//				ResourceName: resourceName,
//			},
//			{
//				Config: updateConfig,
//				Check: resource.ComposeTestCheckFunc(
//					testAccountExists(resourceName),
//					resource.TestCheckResourceAttr(resourceName, "description", description),
//					resource.TestCheckResourceAttrSet(resourceName, "id"),
//					resource.TestCheckResourceAttr(resourceName, "name", name),
//					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
//					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
//					resource.TestCheckResourceAttr(resourceName, "execution_subject_keys.0", executionKeys[0]),
//					resource.TestCheckResourceAttr(resourceName, "audience", updatedAudience),
//				),
//				ResourceName: resourceName,
//			},
//		},
//	})
//}
//
//func testGenericOidcAccountBasic(localName string, name string, description string, tenantedDeploymentParticipation core.TenantedDeploymentMode, execution_subject_keys []string, audience string) string {
//
//	execKeysStr := fmt.Sprintf(`["%s"]`, strings.Join(execution_subject_keys, `", "`))
//
//	return fmt.Sprintf(`resource "octopusdeploy_generic_oidc_account" "%s" {
//		description = "%s"
//		name = "%s"
//		tenanted_deployment_participation = "%s"
//		execution_subject_keys = %s
//		audience = "%s"
//	}
//
//	data "octopusdeploy_accounts" "test" {
//		ids = [octopusdeploy_generic_oidc_account.%s.id]
//	}`, localName, description, name, tenantedDeploymentParticipation, execKeysStr, audience, localName)
//}
