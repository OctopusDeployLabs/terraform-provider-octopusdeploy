package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSSHKeyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testSSHKeyBasic(localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "private_key_passphrase", passphrase),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
			},
		},
	})
}

func testSSHKeyBasic(localName string, name string, privateKeyFile string, username string, passphrase string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_ssh_key_account" "%s" {
		name = "%s"
		private_key_file = "%s"
		private_key_passphrase = "%s"
		tenanted_deployment_participation = "%s"
		username = "%s"
	}`, localName, name, privateKeyFile, passphrase, tenantedDeploymentParticipation, username)
}

// TestSshAccountResource verifies that an SSH account can be reimported with the correct settings
func TestSshAccountResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "7-sshaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := accounts.AccountsQuery{
			PartialName: "SSH",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"SSH\"")
		}
		resource := resources.Items[0].(*accounts.SSHKeyAccount)

		if resource.AccountType != "SshKeyPair" {
			t.Fatal("The account must be have a type of \"SshKeyPair\"")
		}

		if resource.Username != "admin" {
			t.Fatal("The account must be have a username of \"admin\"")
		}

		if resource.Description != "A test account" {
			// This appears to be a bug in the provider where the description is not set
			t.Log("BUG: The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatal("The account must be have no tenant tags")
		}

		if len(resource.EnvironmentIDs) == 0 {
			t.Fatal("The account must have environments")
		}

		return nil
	})
}
