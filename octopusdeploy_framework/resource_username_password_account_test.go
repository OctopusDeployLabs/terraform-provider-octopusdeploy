package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccUsernamePasswordBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_username_password_account." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := testUsernamePasswordBasic(localName, description, name, username, password, tenantedDeploymentParticipation)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "password"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
				ResourceName: resourceName,
			},
		},
	})
}

func testAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Accounts.GetByID(accountID); err != nil {
			return err
		}

		return nil
	}
}

func testUsernamePasswordBasic(localName string, description string, name string, username string, password string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_username_password_account" "%s" {
		description                       = "%s"
		name                              = "%s"
		password                          = "%s"
		tenanted_deployment_participation = "%s"
		username                          = "%s"
	}`, localName, description, name, password, tenantedDeploymentParticipation, username)
}

// TestUsernamePasswordVariableResource verifies that a project variable referencing a username/password account
// can be created
func TestUsernamePasswordVariableResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "54-usernamepasswordvariable", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := projects.Get(client, newSpaceId, query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	resource := resources.Items[0]

	projectVariables, err := variables.GetVariableSet(client, newSpaceId, resource.VariableSetID)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(projectVariables.Variables) != 1 {
		t.Fatalf("The project must have 1 variable.")
	}

	if projectVariables.Variables[0].Name != "UsernamePasswordVariable" {
		t.Fatalf("The variable must be called UsernamePasswordVariable.")
	}

	if projectVariables.Variables[0].Type != "UsernamePasswordAccount" {
		t.Fatalf("The variable must have type of UsernamePasswordAccount.")
	}
}
