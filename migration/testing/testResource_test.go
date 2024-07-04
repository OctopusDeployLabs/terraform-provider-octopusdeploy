package migrationTesting

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"v5": func() (tfprotov5.ProviderServer, error) {
			provider := octopusdeploy.Provider().GRPCProvider()
			return provider, nil
		},
	}
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"v6": func() (tfprotov6.ProviderServer, error) {
			provider := providerserver.NewProtocol6(octopusdeployv6.NewOctopusDeployProviderV6())()
			return provider, nil
		},
	}
}

func TestResource_UpgradeFromVersion(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {

		formattedConfig := fmt.Sprintf(`provider "octopusdeploy" {
  address = "%s"
  api_key = "%s"
}
resource "octopusdeploy_project_group" "MigrationProjects" {
  							name        = "Migration Projects"
  							description = "My Migration Projects group"
						}
`, container.URI, test.ApiKey)

		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"octopusdeploy": {
							VersionConstraint: "0.21.1",
							Source:            "OctopusDeployLabs/octopusdeploy",
						},
					},
					Config: formattedConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("octopusdeploy_project_group.MigrationProjects", "name", "Migration Projects"),
						/* ... */
					),
				},
				{
					ProtoV6ProviderFactories: protoV6ProviderFactories(),
					Config: `resource "octopusdeploy_project_group" "MigrationProjects" {
							name        = "Migration Projects"
							description = "My Migration Projects group"
						}`,
					// ConfigPlanChecks is a terraform-plugin-testing feature.
					// If acceptance testing is still using terraform-plugin-sdk/v2,
					// use `PlanOnly: true` instead. When migrating to
					// terraform-plugin-testing, switch to `ConfigPlanChecks` or you
					// will likely experience test failures.
					PlanOnly: true,
					//ConfigPlanChecks: resource.ConfigPlanChecks{
					//	PreApply: []plancheck.PlanCheck{
					//		plancheck.ExpectEmptyPlan(),
					//	},
					//},
				},
			},
		})
		return nil
	})
}

func Test_Existing(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		terraformTest := &terraform.Options{
			TerraformDir: "../examples/Project-Group-Creation",
		}

		defer terraform.Destroy(t, terraformTest)

		if _, err := terraform.InitE(t, terraformTest); err != nil {
			fmt.Println(err)
		}

		if _, err := terraform.PlanE(t, terraformTest); err != nil {
			fmt.Println(err)
		}

		if _, err := terraform.ApplyE(t, terraformTest); err != nil {
			fmt.Println(err)
		}
		return nil
	})
}

func TestDataSource_UpgradeFromVersion(t *testing.T) {
	/* ... */
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.21.1",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: `data "provider_datasource" "test" {
                            /* ... */
                        }
                        
                        resource "terraform_data" "test" {
                            input = data.provider_datasource.test
                        }`,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config: `data "provider_datasource" "test" {
                            /* ... */
                        }
                        
                        resource "terraform_data" "test" {
                            input = data.provider_datasource.test
                        }`,
				// ConfigPlanChecks is a terraform-plugin-testing feature.
				// If acceptance testing is still using terraform-plugin-sdk/v2,
				// use `PlanOnly: true` instead. When migrating to
				// terraform-plugin-testing, switch to `ConfigPlanChecks` or you
				// will likely experience test failures.
				PlanOnly: true,
				//ConfigPlanChecks: resource.ConfigPlanChecks{
				//	PreApply: []plancheck.PlanCheck{
				//		plancheck.ExpectEmptyPlan(),
				//	},
				//},
			},
		},
	})
}
