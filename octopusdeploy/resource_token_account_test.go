package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTokenAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_token_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted
	token := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccountCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testTokenAccountBasic(localName, description, name, tenantedDeploymentParticipation, token),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(resourceName, "token", token),
				),
			},
		},
	})
}

func testTokenAccountBasic(localName string, description string, name string, tenantedDeploymentParticipation core.TenantedDeploymentMode, token string) string {
	return fmt.Sprintf(`resource "octopusdeploy_token_account" "%s" {
		description                       = "%s"
		name                              = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = []
		token                             = "%s"
	}`, localName, description, name, tenantedDeploymentParticipation, token)
}

// TestTokenAccountResource verifies that a token account can be reimported with the correct settings
func TestTokenAccountResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "9-tokenaccount", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := accounts.AccountsQuery{
		PartialName: "Token",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Accounts.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an account called \"Token\"")
	}
	resource := resources.Items[0].(*accounts.TokenAccount)

	if resource.AccountType != "Token" {
		t.Fatal("The account must be have a type of \"Token\"")
	}

	if !resource.Token.HasValue {
		t.Fatal("The account must be have a token")
	}

	if resource.Description != "A test account" {
		t.Fatal("The account must be have a description of \"A test account\"")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The account must be have a tenanted deployment participation of \"Untenanted\"")
	}

	if len(resource.TenantTags) != 0 {
		t.Fatal("The account must be have no tenant tags")
	}
}
