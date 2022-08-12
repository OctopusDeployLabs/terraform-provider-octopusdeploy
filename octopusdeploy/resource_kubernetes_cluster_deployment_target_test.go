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
		CheckDestroy: testDeploymentTargetCheckDestroy,
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

func TestAccKubernetesClusterDeploymentTargetAws(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountAccessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountSecretKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetAws(
					accountLocalName,
					accountName,
					accountAccessKey,
					accountSecretKey,
					environmentLocalName,
					environmentName,
					userRoleLocalName,
					userRoleName,
					localName,
					name,
					clusterURL,
					clusterName),
			},
		},
	})
}

func TestAccKubernetesClusterDeploymentTargetGcp(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	project := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	region := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetGcp(
					accountLocalName,
					accountName,
					accountUsername,
					environmentLocalName,
					environmentName,
					userRoleLocalName,
					userRoleName,
					localName,
					name,
					clusterURL,
					clusterName,
					project,
					region),
			},
		},
	})
}

func testAccKubernetesClusterDeploymentTargetBasic(accountLocalName string, accountName string, accountUsername string, environmentLocalName string, environmentName string, userRoleLocalName string, userRoleName string, localName string, name string, clusterURL string) string {
	usernamePasswordAccountID := "${octopusdeploy_username_password_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	return fmt.Sprintf(testUsernamePasswordMinimum(accountLocalName, accountName, accountUsername)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   authentication {
		     account_id = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, usernamePasswordAccountID)
}

func testAccKubernetesClusterDeploymentTargetGcp(
	accountLocalName string,
	accountName string,
	accountJSONKey string,
	environmentLocalName string,
	environmentName string,
	userRoleLocalName string,
	userRoleName string,
	localName string,
	name string,
	clusterURL string,
	clusterName string,
	project string,
	region string) string {
	gcpAccountID := "${octopusdeploy_gcp_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	return fmt.Sprintf(testGcpAccount(accountLocalName, accountName, accountJSONKey)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   gcp_account_authentication {
		     account_id = "%s"
			 cluster_name = "%s"
			 project = "%s"
			 region = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, gcpAccountID, clusterName, project, region)
}

func testAccKubernetesClusterDeploymentTargetAws(
	accountLocalName string,
	accountName string,
	accountAccessKey string,
	accountSecretKey string,
	environmentLocalName string,
	environmentName string,
	userRoleLocalName string,
	userRoleName string,
	localName string,
	name string,
	clusterURL string,
	clusterName string) string {
	awsAccountID := "${octopusdeploy_aws_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	return fmt.Sprintf(testAwsAccount(accountLocalName, accountName, accountAccessKey, accountSecretKey)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   aws_account_authentication {
		     account_id = "%s"
			 cluster_name = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, awsAccountID, clusterName)
}
