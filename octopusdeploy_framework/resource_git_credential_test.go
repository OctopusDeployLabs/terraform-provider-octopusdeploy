package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"path/filepath"
	"testing"
)

func TestGitCredentialBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_git_credential." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
				),
				Config: testGitCredential(localName, name, description),
			},
		},
	})
}

func testGitCredential(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_git_credential" "%s" {
		name = "%s"
		  description  = "%s"
		  username     = "git_user"
		  password     = "secret_password"
	}`, localName, name, description)
}

// TestGitCredentialsResource verifies that a git credential can be reimported with the correct settings
func TestGitCredentialsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "22-gitcredentialtest", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "22a-gitcredentialtestds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "22a-gitcredentialtestds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup == "" {
		t.Fatal("The target lookup did not succeed.")
	}
}
