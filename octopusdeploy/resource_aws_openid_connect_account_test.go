package octopusdeploy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSOIDCAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_aws_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	roleArn := "arn:aws:iam::sourceAccountId:roleroleName"
	sessionDuration := "3600"
	executionKeys := []string{"space"}
	healthKeys := []string{"target"}
	accountKeys := []string{"type"}

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
					resource.TestCheckResourceAttr(prefix, "role_arn", roleArn),
					resource.TestCheckResourceAttr(prefix, "session_duration", sessionDuration),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys", healthKeys[0]),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys", accountKeys[0]),
				),
				Config: testAwsOIDCAccountBasic(localName, name, description, roleArn, sessionDuration, tenantedDeploymentParticipation, executionKeys, healthKeys, accountKeys),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "role_arn", roleArn),
					resource.TestCheckResourceAttr(prefix, "session_duration", sessionDuration),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys", executionKeys[0]),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys", healthKeys[0]),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys", accountKeys[0]),
				),
				Config: testAwsOIDCAccountBasic(localName, name, description, roleArn, sessionDuration, tenantedDeploymentParticipation, executionKeys, healthKeys, accountKeys),
			},
		},
	})
}

func testAwsOIDCAccountBasic(localName string, name string, description string, roleArn string, sessionDuration string, tenantedDeploymentParticipation core.TenantedDeploymentMode, execution_subject_keys []string, health_subject_keys []string, account_test_subject_keys []string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_openid_connect_account" "%s" {
		description = "%s"
		name = "%s"
		role_arn = "%s"
		tenanted_deployment_participation = "%s"
		execution_subject_keys = %s
		health_subject_keys = %s
		account_test_subject_keys = %s
		session_duration = %s
	}
	
	data "octopusdeploy_accounts" "test" {
		ids = [octopusdeploy_aws_openid_connect_account.%s.id]
	}`, localName, description, name, roleArn, tenantedDeploymentParticipation, StringArrayToTerraformArrayFormat(execution_subject_keys), StringArrayToTerraformArrayFormat(health_subject_keys), StringArrayToTerraformArrayFormat(account_test_subject_keys), sessionDuration, localName)
}

func testAwsOIDCAccount(localName string, name string, roleArn string, sessionDuration string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_openid_connect_account" "%s" {
		name       = "%s"
		role_arn   = "%s"
		session_duration = %s
	}`, localName, name, roleArn, sessionDuration)
}

func StringArrayToTerraformArrayFormat(slice []string) string {
	// Convert each element in the slice to a quoted string
	for i, v := range slice {
		slice[i] = fmt.Sprintf("%q", v)
	}
	// Join the quoted elements with commas
	joined := strings.Join(slice, ", ")
	// Format the string to include square brackets
	formatted := fmt.Sprintf("[%s]", joined)
	return formatted
}
