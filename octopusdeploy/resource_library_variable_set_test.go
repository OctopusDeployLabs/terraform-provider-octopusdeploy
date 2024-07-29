package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (suite *IntegrationTestSuite) TestAccOctopusDeployLibraryVariableSetBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_library_variable_set." + localName
	t := suite.T()

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testLibraryVariableSetDestroy,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testLibraryVariableSetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
			},
		},
	})
}

func (suite *IntegrationTestSuite) TestAccOctopusDeployLibraryVariableSetComplex() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_library_variable_set." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy:             testLibraryVariableSetDestroy,
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testLibraryVariableSetBasic(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "template.#", "0"),
				),
			},
			{
				Config: testLibraryVariableSetBasicWithDescription(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "template.#", "0"),
				),
			},
			{
				Config: testLibraryVariableSetComplex(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "template.#", "3"),
					resource.TestCheckResourceAttr(prefix, "template.1.default_value", "wjehqwjkehwqkejhqwe"),
				),
			},
		},
	})
}

func (suite *IntegrationTestSuite) TestAccOctopusDeployLibraryVariableSetWithUpdate() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_library_variable_set." + localName

	dataLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	dataPrefix := "data.octopusdeploy_library_variable_sets." + dataLocalName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy:             testLibraryVariableSetDestroy,
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// create variable set with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testLibraryVariableSetBasic(localName, name),
			},
			// create update it with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testLibraryVariableSetBasicWithDescription(localName, name, description),
			},
			// update again by remove its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOctopusDeployLibraryVariableSetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", ""),
				),
				Config: testLibraryVariableSetBasic(localName, name),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibraryVariableSetDataSourceID(dataPrefix),
				),
				Config: testLibraryVariableSetsData(dataLocalName, name),
			},
		},
	})
}

func testLibraryVariableSetBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_library_variable_set" "%s" {
		name = "%s"
	}`, localName, name)
}

func testLibraryVariableSetBasicWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_library_variable_set" "%s" {
		name        = "%s"
		description = "%s"
	}`, localName, name, description)
}

func testLibraryVariableSetsData(localName string, name string) string {
	return fmt.Sprintf(`data "octopusdeploy_library_variable_sets" "%s" {
		partial_name = "%s"
	}`, localName, name)
}

func testLibraryVariableSetComplex(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_library_variable_set" "%s" {
		description = "This is the description."
		name        = "%s"

		template {
			default_value    = "Default Value???"
			display_settings = {
				"Octopus.ControlType" = "SingleLineText"
			}
			help_text        = "This is the help text"
			label            = "Test Label"
			name             = "Test Template"
		}

		template {
			default_value    = "wjehqwjkehwqkejhqwe"
			display_settings = {
				"Octopus.ControlType" = "MultiLineText"
			}
			help_text        = "jhasdkjashdaksjhd"
			label            = "alsdjhaldksh"
			name             = "Another Variable???"
		}

		template {
			default_value    = "qweq|qwe"
			display_settings = {
				"Octopus.ControlType" = "MultiLineText"
			}
			help_text        = "qwe"
			label            = "qwe"
			name             = "weq"
		}
	}`, localName, name)
}

func testLibraryVariableSetDestroy(s *terraform.State) error {
	if err := destroyHelperLibraryVariableSet(s); err != nil {
		return err
	}
	return nil
}

func testAccCheckOctopusDeployLibraryVariableSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := existsHelperLibraryVariableSet(s, octoClient); err != nil {
			return err
		}
		return nil
	}
}

func destroyHelperLibraryVariableSet(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		libraryVariableSetID := rs.Primary.ID
		libraryVariableSet, err := octoClient.LibraryVariableSets.GetByID(libraryVariableSetID)
		if err == nil {
			if libraryVariableSet != nil {
				return fmt.Errorf("library variable set (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func existsHelperLibraryVariableSet(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_library_variable_set" {
			if _, err := client.LibraryVariableSets.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving library variable set %s", err)
			}
		}
	}
	return nil
}

func testAccCheckLibraryVariableSetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find library variable set data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot library variable set source ID not set")
		}
		return nil
	}
}
