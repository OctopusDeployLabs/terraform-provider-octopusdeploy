package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccOctopusDeployServiceAccountOIDCIdentity(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_service_account_oidc_identity." + localName

	localUserName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userPrefix := " octopusdeploy_user." + localUserName

	userData := users.User{
		DisplayName:  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		EmailAddress: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + "@test.com",
		Username:     acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}

	data := serviceaccounts.OIDCIdentity{
		Name:             acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ServiceAccountID: userPrefix + ".id",
		Issuer:           "https://token.actions.githubusercontent.com",
		Subject:          "repo:test/test:environment:test",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testServiceAccountIdentityConfig(localName, localUserName, data, userData),
				Check: resource.ComposeTestCheckFunc(
					testScriptModuleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", data.Name),
					resource.TestCheckResourceAttr(prefix, "issuer", data.Issuer),
					resource.TestCheckResourceAttr(prefix, "subject", data.Subject),
				),
			},
			{
				Config: testServiceAccountIdentityUpdate(localName, localUserName, data, userData),
				Check: resource.ComposeTestCheckFunc(
					testScriptModuleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", data.Name+"-updated"),
					resource.TestCheckResourceAttr(prefix, "issuer", data.Issuer),
					resource.TestCheckResourceAttr(prefix, "subject", data.Subject),
				),
			},
		},
	})
}

func testServiceAccountIdentityConfig(localName string, localUserName string, data serviceaccounts.OIDCIdentity, userData users.User) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user" "%s" {
		display_name = "%s"
		email_address = "%s"
		is_active = true
		is_service = true
		username = "%s"
	}
	resource "octopusdeploy_service_account_oidc_identity" "%s" {
		name        = "%s"
		service_account_id = %s
		issuer = "%s"
		subject = "%s"
	}`,
		localUserName,
		userData.DisplayName,
		userData.EmailAddress,
		userData.Username,
		localName,
		data.Name,
		data.ServiceAccountID,
		data.Issuer,
		data.Subject)
}

func testServiceAccountIdentityUpdate(localName string, localUserName string, data serviceaccounts.OIDCIdentity, userData users.User) string {
	data.Name = data.Name + "-updated"
	return testServiceAccountIdentityUpdate(localName, localUserName, data, userData)
}
