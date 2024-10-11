package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"path/filepath"
	"strconv"
	"testing"
)

func TestAccUserImportBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserBasic(localName, displayName, true, false, password, username, emailAddress),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUserBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user." + localName

	isActive := true
	isService := false
	displayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	emailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "email_address", emailAddress),
					resource.TestCheckResourceAttr(resourceName, "is_active", strconv.FormatBool(isActive)),
					resource.TestCheckResourceAttr(resourceName, "is_service", strconv.FormatBool(isService)),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
				Config: testAccUserBasic(localName, displayName, isActive, isService, password, username, emailAddress),
			},
			{
				//Config:                  testAccUserImport(localName, username),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateIdFunc:       testAccUserImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccUserImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccUserImport(localName string, username string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {}`, localName)
}

func testAccUserBasic(localName string, displayName string, isActive bool, isService bool, password string, username string, emailAddress string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {
		display_name  = "%s"
		email_address = "%s"
		is_active     = %v
		is_service    = %v
		password      = "%s"
		username      = "%s"

		identity {
			provider = "Octopus ID"
			claim {
				name = "email"
				is_identifying_claim = true
				value = "%s"
			}
			claim {
				name = "dn"
				is_identifying_claim = false
				value = "%s"
			}
		}
	}`, localName, displayName, emailAddress, isActive, isService, password, username, emailAddress, displayName)
}

func testUserExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		userID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := users.GetByID(octoClient, userID); err != nil {
			return err
		}

		return nil
	}
}

func testAccUserCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_user" {
			continue
		}

		_, err := users.GetByID(octoClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("user (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestProjectTerraformPackageScriptExport verifies that users and teams can be reimported
func TestUsersAndTeams(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "43-users", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "43a-usersds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)

	if err != nil {
		t.Fatal(err.Error())
	}

	err = func() error {
		query := users.UsersQuery{
			Filter: "Service Account",
			IDs:    nil,
			Skip:   0,
			Take:   1,
		}

		resources, err := client.Users.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a user called \"Service Account\"")
		}

		resource := resources.Items[0]

		if resource.Username != "saccount" {
			t.Fatalf("Account must have a username \"saccount\"")
		}

		if resource.EmailAddress != "a@a.com" {
			t.Fatalf("Account must have a email \"a@a.com\"")
		}

		if !resource.IsService {
			t.Fatalf("Account must be a service account")
		}

		if !resource.IsActive {
			t.Fatalf("Account must be active")
		}

		return nil
	}()

	if err != nil {
		t.Fatal(err.Error())
	}

	err = func() error {
		query := users.UsersQuery{
			Filter: "Bob Smith",
			IDs:    nil,
			Skip:   0,
			Take:   1,
		}

		resources, err := users.Get(client, "", query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a user called \"Service Account\"")
		}

		resource := resources.Items[0]

		if resource.Username != "bsmith" {
			t.Fatalf("Regular account must have a username \"bsmith\"")
		}

		if resource.EmailAddress != "bob.smith@example.com" {
			t.Fatalf("Regular account must have a email \"bob.smith@example.com\"")
		}

		if resource.IsService {
			t.Fatalf("Account must not be a service account")
		}

		if resource.IsActive {
			t.Log("BUG: Account must not be active")
		}

		return nil
	}()

	if err != nil {
		t.Fatal(err.Error())
	}

	err = func() error {
		query := teams.TeamsQuery{
			IDs:           nil,
			IncludeSystem: false,
			PartialName:   "Deployers",
			Skip:          0,
			Spaces:        nil,
			Take:          1,
		}

		resources, err := client.Teams.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a team called \"Deployers\"")
		}

		resource := resources.Items[0]

		if len(resource.MemberUserIDs) != 1 {
			t.Fatalf("Team must have one user")
		}

		return nil
	}()

	if err != nil {
		t.Fatal(err.Error())
	}

	// Verify the environment data lookups work
	teams, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "43a-usersds"), "teams_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if teams == "" {
		t.Fatal("The teams lookup failed.")
	}

	roles, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "43a-usersds"), "roles_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if roles == "" {
		t.Fatal("The roles lookup failed.")
	}

	users, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "43a-usersds"), "users_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if users == "" {
		t.Fatal("The users lookup failed.")
	}
}
