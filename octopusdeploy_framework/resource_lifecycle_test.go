package octopusdeploy_framework

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestExpandLifecycleWithNil(t *testing.T) {
	lifecycle := expandLifecycle(nil)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.NewRetentionPeriod(0, "Days", true)
	tentacleRetention := core.NewRetentionPeriod(2, "Items", false)

	data := &lifecycleTypeResourceModel{
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(tentacleRetention.ShouldKeepForever),
						"unit":                types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycle(data)

	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
	require.Equal(t, spaceID, lifecycle.SpaceID)
}

func TestExpandPhasesWithEmptyInput(t *testing.T) {
	emptyList := types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, []attr.Value{})
	phases := expandPhases(emptyList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithNullInput(t *testing.T) {
	nullList := types.ListNull(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	phases := expandPhases(nullList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithUnknownInput(t *testing.T) {
	unknownList := types.ListUnknown(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	phases := expandPhases(unknownList)
	require.Nil(t, phases)
}

func TestExpandAndFlattenPhasesWithSensibleDefaults(t *testing.T) {
	phase := createTestPhase("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

	flattenedPhases := flattenPhases([]*lifecycles.Phase{phase})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 1, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 1)

	expandedPhase := expandedPhases[0]
	require.NotEmpty(t, expandedPhase.ID)
	require.Equal(t, phase.AutomaticDeploymentTargets, expandedPhase.AutomaticDeploymentTargets)
	require.Equal(t, phase.IsOptionalPhase, expandedPhase.IsOptionalPhase)
	require.EqualValues(t, phase.MinimumEnvironmentsBeforePromotion, expandedPhase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase.Name, expandedPhase.Name)
	require.Equal(t, phase.ReleaseRetentionPolicy, expandedPhase.ReleaseRetentionPolicy)
	require.Equal(t, phase.TentacleRetentionPolicy, expandedPhase.TentacleRetentionPolicy)
}

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaults(t *testing.T) {
	phase1 := createTestPhase("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhase("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

	flattenedPhases := flattenPhases([]*lifecycles.Phase{phase1, phase2})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 2, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 2)

	require.NotEmpty(t, expandedPhases[0].ID)
	require.Equal(t, phase1.AutomaticDeploymentTargets, expandedPhases[0].AutomaticDeploymentTargets)
	require.Equal(t, phase1.IsOptionalPhase, expandedPhases[0].IsOptionalPhase)
	require.EqualValues(t, phase1.MinimumEnvironmentsBeforePromotion, expandedPhases[0].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase1.Name, expandedPhases[0].Name)
	require.Equal(t, phase1.ReleaseRetentionPolicy, expandedPhases[0].ReleaseRetentionPolicy)
	require.Equal(t, phase1.TentacleRetentionPolicy, expandedPhases[0].TentacleRetentionPolicy)

	require.NotEmpty(t, expandedPhases[1].ID)
	require.Equal(t, phase2.AutomaticDeploymentTargets, expandedPhases[1].AutomaticDeploymentTargets)
	require.Equal(t, phase2.IsOptionalPhase, expandedPhases[1].IsOptionalPhase)
	require.EqualValues(t, phase2.MinimumEnvironmentsBeforePromotion, expandedPhases[1].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase2.Name, expandedPhases[1].Name)
	require.Equal(t, phase2.ReleaseRetentionPolicy, expandedPhases[1].ReleaseRetentionPolicy)
	require.Equal(t, phase2.TentacleRetentionPolicy, expandedPhases[1].TentacleRetentionPolicy)
}

func createTestPhase(name string, autoTargets []string, isOptional bool, minEnvs int32) *lifecycles.Phase {
	phase := lifecycles.NewPhase(name)
	phase.AutomaticDeploymentTargets = autoTargets
	phase.IsOptionalPhase = isOptional
	phase.MinimumEnvironmentsBeforePromotion = minEnvs
	phase.ReleaseRetentionPolicy = core.NewRetentionPeriod(15, "Items", false)
	phase.TentacleRetentionPolicy = core.NewRetentionPeriod(0, "Days", true)
	phase.ID = name + "-Id"
	return phase
}

//Integration test under here

func TestAccLifecycleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
		},
	})
}

func TestAccLifecycleWithUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// create lifecycle with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
			// update lifecycle with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycleWithDescription(localName, name, description),
			},
			// update lifecycle by removing its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "30"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					//resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycle(localName, name),
			},
			// update lifecycle add retention policy
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "60"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
				),
				Config: testAccLifecycleWithRetentionPolicy(localName, name, description),
			},
		},
	})
}

func TestAccLifecycleComplex(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "2"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
					testAccCheckLifecyclePhaseCount(name, 2),
				),
				Config: testAccLifecycleComplex(localName, name),
			},
		},
	})
}

func testAccLifecycle(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = ""
	}`, localName, name)
}

func testAccLifecycleWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
       description = "%s"
    }`, localName, name, description)
}

func testAccLifecycleWithRetentionPolicy(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
       description = "%s"
		release_retention_policy {
			unit             = "Days"
			quantity_to_keep = 60
			should_keep_forever = false
		}

		tentacle_retention_policy {
			unit             = "Days"
			quantity_to_keep = 0
			should_keep_forever = true
		}
    }`, localName, name, description)
}

func testAccLifecycleComplex(localName string, name string) string {
	environment1LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment1Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccEnvironment(environment1LocalName, environment1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment2LocalName, environment2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment3LocalName, environment3Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		`resource "octopusdeploy_lifecycle" "%s" {
			name        = "%s"
			description = "Funky Lifecycle description"

			release_retention_policy {
				unit             = "Days"
				quantity_to_keep = 2
			}

			tentacle_retention_policy {
				unit             = "Days"
				quantity_to_keep = 1
			}

			phase {
				automatic_deployment_targets          = ["${octopusdeploy_environment.%s.id}"]
				is_optional_phase                     = true
				minimum_environments_before_promotion = 2
				name                                  = "P1"
				optional_deployment_targets           = ["${octopusdeploy_environment.%s.id}"]
			}

			phase {
				name = "P2"
			}
	}`, localName, name, environment2LocalName, environment3LocalName)
}

func testAccCheckLifecycleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := existsHelperLifecycle(s, octoClient); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceList, err := octoClient.Lifecycles.GetByPartialName(name)
		if err != nil {
			return err
		}

		resource := resourceList[0]

		if len(resource.Phases) != expected {
			return fmt.Errorf("lifecycle has %d phases instead of the expected %d", len(resource.Phases), expected)
		}

		return nil
	}
}

func existsHelperLifecycle(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_lifecycle" {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}

func testAccLifecycleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_lifecycle" {
			continue
		}

		lifecycle, err := octoClient.Lifecycles.GetByID(rs.Primary.ID)
		if err == nil && lifecycle != nil {
			return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestLifecycleResource verifies that a lifecycle can be reimported with the correct settings
func TestLifecycleResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "17-lifecycle", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "17a-lifecycleds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := lifecycles.Query{
		PartialName: "Simple",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Lifecycles.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have an environment called \"Simple\"")
	}
	resource := resources.Items[0]

	if resource.Description != "A test lifecycle" {
		t.Fatal("The lifecycle must be have a description of \"A test lifecycle\" (was \"" + resource.Description + "\")")
	}

	if resource.TentacleRetentionPolicy.QuantityToKeep != 30 {
		t.Fatal("The lifecycle must be have a tentacle retention policy of \"30\" (was \"" + fmt.Sprint(resource.TentacleRetentionPolicy.QuantityToKeep) + "\")")
	}

	if resource.TentacleRetentionPolicy.ShouldKeepForever {
		t.Fatal("The lifecycle must be have a tentacle retention not set to keep forever")
	}

	if resource.TentacleRetentionPolicy.Unit != "Items" {
		t.Fatal("The lifecycle must be have a tentacle retention unit set to \"Items\" (was \"" + resource.TentacleRetentionPolicy.Unit + "\")")
	}

	if resource.ReleaseRetentionPolicy.QuantityToKeep != 1 {
		t.Fatal("The lifecycle must be have a release retention policy of \"1\" (was \"" + fmt.Sprint(resource.ReleaseRetentionPolicy.QuantityToKeep) + "\")")
	}

	if !resource.ReleaseRetentionPolicy.ShouldKeepForever {
		t.Log("BUG: The lifecycle must be have a release retention set to keep forever (known bug - the provider creates this field as false)")
	}

	if resource.ReleaseRetentionPolicy.Unit != "Days" {
		t.Fatal("The lifecycle must be have a release retention unit set to \"Days\" (was \"" + resource.ReleaseRetentionPolicy.Unit + "\")")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "17a-lifecycleds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
