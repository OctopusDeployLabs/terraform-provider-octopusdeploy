package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesClusterDeploymentTargetBasic(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	newClusterURL := "http://www.example.com"

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetBasic(accountLocalName, accountName, accountUsername, environmentLocalName, environmentName, userRoleLocalName, userRoleName, localName, name, clusterURL),
			},
			{
				Config: testAccKubernetesClusterDeploymentTargetBasic(accountLocalName, accountName, accountUsername, environmentLocalName, environmentName, userRoleLocalName, userRoleName, localName, name, newClusterURL),
			},
		},
	})
}

func testAccKubernetesClusterDeploymentTargetBasic(accountLocalName string, accountName string, accountUsername string, environmentLocalName string, environmentName string, userRoleLocalName string, userRoleName string, localName string, name string, clusterURL string) string {
	usernamePasswordAccountID := "${octopusdeploy_username_password_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	return fmt.Sprintf(testUsernamePasswordMinimum(accountLocalName, accountName, accountUsername)+"\n"+
		testEnvironmentMinimum(environmentLocalName, environmentName)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]

		   authentication {
		     account_id = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, usernamePasswordAccountID)
}
