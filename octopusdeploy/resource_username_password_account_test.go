package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func (suite *IntegrationTestSuite) TestAccUsernamePasswordBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_username_password_account." + localName
	t := suite.T()

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	config := testUsernamePasswordBasic(localName, description, name, username, password, tenantedDeploymentParticipation)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccountCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
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

func testUsernamePasswordBasic(localName string, description string, name string, username string, password string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_username_password_account" "%s" {
		description                       = "%s"
		name                              = "%s"
		password                          = "%s"
		tenanted_deployment_participation = "%s"
		username                          = "%s"
	}`, localName, description, name, password, tenantedDeploymentParticipation, username)
}

func testUsernamePasswordMinimum(localName string, name string, username string) string {
	return fmt.Sprintf(`resource "octopusdeploy_username_password_account" "%s" {
		name     = "%s"
		username = "%s"
	}`, localName, name, username)
}

// TestUsernamePasswordVariableResource verifies that a project variable referencing a username/password account
// can be created
func (suite *IntegrationTestSuite) TestUsernamePasswordVariableResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()
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
