package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccMapProcessStepFromStateWithAllAttributes(t *testing.T) {
	ctx := context.Background()
	primaryPackage := &schemas.ProcessStepPackageReferenceResourceModel{
		PackageID:           types.StringValue("Packages-0"),
		FeedID:              types.StringValue("Feeds-0"),
		AcquisitionLocation: types.StringValue("ExecutionTarget"),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Extract": types.StringValue("True"),
		}),
	}
	primaryPackage.ID = types.StringValue("00000000-0000-0000-0000-000000000022")
	state := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Step One"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Always"),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.MaxParallelism": types.StringValue("2"),
			"Octopus.Action.TargetRoles":    types.StringValue("agent-1,agent-2"),
		}),
		Type:               types.StringValue("Octopus.Script"),
		Slug:               types.StringValue("step-one"),
		IsRequired:         types.BoolValue(true),
		IsDisabled:         types.BoolValue(false),
		Notes:              types.StringValue(`Some notes`),
		WorkerPoolID:       types.StringValue("WorkerPools-1"),
		WorkerPoolVariable: types.StringValue("#{Environment.WorkerPools.Default}"),
		TenantTags: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("tag-1"),
			types.StringValue("tag-2"),
		}),
		Environments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-1"),
			types.StringValue("Environments-2"),
		}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-13"),
		}),
		Channels: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Channels-1"),
		}),
		Container: &schemas.ProcessStepActionContainerModel{
			FeedID: types.StringValue("Feeds-1"),
			Image:  types.StringValue("docker.io/library/dummy:latest"),
		},
		GitDependencies: types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), map[string]attr.Value{
			"script-folder": types.ObjectValueMust(
				schemas.ProcessStepGitDependencyAttributeTypes(),
				map[string]attr.Value{
					"repository_uri":      types.StringValue("git://test.repository.fi"),
					"default_branch":      types.StringValue("main"),
					"git_credential_type": types.StringValue("UsernamePassword"),
					"file_path_filters": types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("directory-a"),
					}),
					"git_credential_id": types.StringValue("GitCredentials-1"),
				},
			),
		}),
		PrimaryPackage: primaryPackage,
		Packages: types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), map[string]attr.Value{
			"script-package": types.ObjectValueMust(
				schemas.ProcessStepPackageReferenceAttributeTypes(),
				map[string]attr.Value{
					"id":                   types.StringValue("00000000-0000-0000-0000-000000000001"),
					"package_id":           types.StringValue("Package-1"),
					"feed_id":              types.StringValue("Feeds-2"),
					"acquisition_location": types.StringValue("#{LocationVariable}"),
					"properties": types.MapValueMust(types.StringType, map[string]attr.Value{
						"Octopus.Package.IsPrimary": types.StringValue("True"),
					}),
				},
			),
		}),
		ExecutionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.RunOnServer":       types.StringValue("True"),
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Step 1, Action 1\""),
		}),
	}

	step := deployments.NewDeploymentStep("Should be replaced with name from the state")

	diags := mapProcessStepFromState(ctx, &state, step)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedStep := &deployments.DeploymentStep{
		Name:               "Step One",
		StartTrigger:       "StartAfterPrevious",
		PackageRequirement: "LetOctopusDecide",
		Condition:          "Always",
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.MaxParallelism": core.NewPropertyValue("2", false),
			"Octopus.Action.TargetRoles":    core.NewPropertyValue("agent-1,agent-2", false),
		},
		Actions: []*deployments.DeploymentAction{
			{
				Name:                 "Step One",
				Slug:                 "step-one",
				ActionType:           "Octopus.Script",
				IsRequired:           true,
				IsDisabled:           false,
				Notes:                "Some notes",
				WorkerPool:           "WorkerPools-1",
				WorkerPoolVariable:   "#{Environment.WorkerPools.Default}",
				TenantTags:           []string{"tag-1", "tag-2"},
				Environments:         []string{"Environments-1", "Environments-2"},
				ExcludedEnvironments: []string{"Environments-13"},
				Channels:             []string{"Channels-1"},
				Container: &deployments.DeploymentActionContainer{
					FeedID: "Feeds-1",
					Image:  "docker.io/library/dummy:latest",
				},
				GitDependencies: []*gitdependencies.GitDependency{
					{
						Name:              "script-folder",
						RepositoryUri:     "git://test.repository.fi",
						DefaultBranch:     "main",
						GitCredentialType: "UsernamePassword",
						FilePathFilters:   []string{"directory-a"},
						GitCredentialId:   "GitCredentials-1",
					},
				},
				Packages: []*packages.PackageReference{
					{
						ID:                  "00000000-0000-0000-0000-000000000022",
						Name:                "", // Empty for primary package
						PackageID:           "Packages-0",
						FeedID:              "Feeds-0",
						AcquisitionLocation: "ExecutionTarget",
						Properties: map[string]string{
							"Extract": "True",
						},
					},
					{
						ID:                  "00000000-0000-0000-0000-000000000001",
						Name:                "script-package",
						PackageID:           "Package-1",
						FeedID:              "Feeds-2",
						AcquisitionLocation: "#{LocationVariable}",
						Properties: map[string]string{
							"Octopus.Package.IsPrimary": "True",
						},
					},
				},
				Properties: map[string]core.PropertyValue{
					"Octopus.Action.RunOnServer":                core.NewPropertyValue("True", false),
					"Octopus.Action.Script.ScriptBody":          core.NewPropertyValue("Write-Host \"Step 1, Action 1\"", false),
					"Octopus.Action.Package.FeedId":             core.NewPropertyValue("Feeds-0", false),
					"Octopus.Action.Package.PackageId":          core.NewPropertyValue("Packages-0", false),
					"Octopus.Action.Package.DownloadOnTentacle": core.NewPropertyValue("True", false),
				},
				Resource: *resources.NewResource(),
			},
		},
		TargetRoles: []string{},
		Resource:    *resources.NewResource(),
	}

	assert.Equal(t, expectedStep, step)
}

func TestAccMapProcessStepFromStateForScriptStep(t *testing.T) {
	ctx := context.Background()
	state := schemas.ProcessStepResourceModel{
		SpaceID:              types.StringValue("Spaces-1"),
		ProcessID:            types.StringValue("Processes-1"),
		Name:                 types.StringValue("Run Script"),
		StartTrigger:         types.StringValue("StartAfterPrevious"),
		PackageRequirement:   types.StringValue("LetOctopusDecide"),
		Condition:            types.StringValue("Success"),
		Type:                 types.StringValue("Octopus.Script"),
		TenantTags:           types.SetValueMust(types.StringType, []attr.Value{}),
		Environments:         types.SetValueMust(types.StringType, []attr.Value{}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{}),
		Channels:             types.SetValueMust(types.StringType, []attr.Value{}),
		Container:            &schemas.ProcessStepActionContainerModel{FeedID: types.StringValue(""), Image: types.StringValue("")},
		GitDependencies:      types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), map[string]attr.Value{}),
		Packages:             types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), map[string]attr.Value{}),
		ExecutionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Minimum attributes\""),
		}),
	}

	step := deployments.NewDeploymentStep("Run Script")

	diags := mapProcessStepFromState(ctx, &state, step)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedStep := &deployments.DeploymentStep{
		Name:               "Run Script",
		StartTrigger:       "StartAfterPrevious",
		PackageRequirement: "LetOctopusDecide",
		Condition:          "Success",
		Properties:         map[string]core.PropertyValue{},
		Actions: []*deployments.DeploymentAction{
			{
				Name:                 "Run Script",
				ActionType:           "Octopus.Script",
				TenantTags:           []string{},
				Environments:         []string{},
				ExcludedEnvironments: []string{},
				Channels:             []string{},
				Container:            &deployments.DeploymentActionContainer{FeedID: "", Image: ""},
				GitDependencies:      []*gitdependencies.GitDependency{},
				Packages:             []*packages.PackageReference{},
				Properties: map[string]core.PropertyValue{
					"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host \"Minimum attributes\"", false),
				},
				Resource: *resources.NewResource(),
			},
		},
		TargetRoles: []string{},
		Resource:    *resources.NewResource(),
	}

	assert.Equal(t, expectedStep, step)
}

func TestAccMapProcessStepToStateWithAllAttributes(t *testing.T) {
	primaryPackage := &packages.PackageReference{
		ID:                  "00000000-0000-0000-0000-000000000101",
		Name:                "",
		PackageID:           "Package-1",
		FeedID:              "Feeds-1",
		AcquisitionLocation: "ExecutionTarget",
		Properties: map[string]string{
			"Octopus.Package.IsPrimary": "True",
		},
	}
	additionalPackage := &packages.PackageReference{
		ID:                  "00000000-0000-0000-0000-000000000102",
		Name:                "unique-name",
		PackageID:           "Package-2",
		FeedID:              "feeds-builtin",
		AcquisitionLocation: "Server",
	}
	gitDependency := &gitdependencies.GitDependency{
		Name:              "this-dependency",
		RepositoryUri:     "git://test.repository.co.nz",
		DefaultBranch:     "default",
		GitCredentialType: "NotSpecified",
		FilePathFilters:   []string{"directory-b"},
		GitCredentialId:   "GitCredential-2",
	}

	action := deployments.NewDeploymentAction("Step One", "Octopus.Script")
	action.SetID("00000000-0000-0000-0000-000000000011")
	action.Name = "Step One"
	action.Slug = "step-one"
	action.ActionType = "Octopus.Script"
	action.IsRequired = true
	action.IsDisabled = false
	action.Notes = "Some notes"
	action.WorkerPool = "WorkerPools-1"
	action.WorkerPoolVariable = "#{Environment.WorkerPools.Default}"
	action.TenantTags = []string{"tag-1", "tag-2"}
	action.Environments = []string{"Environments-1", "Environments-2"}
	action.ExcludedEnvironments = []string{"Environments-13"}
	action.Channels = []string{"Channels-1"}
	action.Container = &deployments.DeploymentActionContainer{
		FeedID: "Feeds-1",
		Image:  "docker.io/library/dummy:latest",
	}
	action.GitDependencies = []*gitdependencies.GitDependency{gitDependency}
	action.Packages = []*packages.PackageReference{primaryPackage, additionalPackage}
	action.Properties = map[string]core.PropertyValue{
		"Octopus.Action.RunOnServer":                core.NewPropertyValue("True", false),
		"Octopus.Action.Script.ScriptBody":          core.NewPropertyValue("Write-Host \"Step 1, Action 1\"", false),
		"Octopus.Action.Package.FeedId":             core.NewPropertyValue("ignored-by-primary-package", false),
		"Octopus.Action.Package.PackageId":          core.NewPropertyValue("ignored-by-primary-package", false),
		"Octopus.Action.Package.DownloadOnTentacle": core.NewPropertyValue("ignored-by-primary-package", false),
	}

	step := deployments.NewDeploymentStep("Step One")
	step.SetID("00000000-0000-0000-0000-000000000001")
	step.StartTrigger = "StartAfterPrevious"
	step.PackageRequirement = "LetOctopusDecide"
	step.Condition = "Success"
	step.Properties = map[string]core.PropertyValue{
		"Octopus.Action.MaxParallelism": core.NewPropertyValue("2", false),
		"Octopus.Action.TargetRoles":    core.NewPropertyValue("agent-1,agent-2", false),
	}
	step.Actions = []*deployments.DeploymentAction{action}

	process := &deployments.DeploymentProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		Steps:     []*deployments.DeploymentStep{step},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
	}
	diags := mapProcessStepToState(deploymentProcessWrapper{process}, step, &state)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedPrimaryPackage := &schemas.ProcessStepPackageReferenceResourceModel{
		PackageID:           types.StringValue(primaryPackage.PackageID),
		FeedID:              types.StringValue(primaryPackage.FeedID),
		AcquisitionLocation: types.StringValue(primaryPackage.AcquisitionLocation),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Package.IsPrimary": types.StringValue("True"),
		}),
	}
	expectedPrimaryPackage.ID = types.StringValue(primaryPackage.ID)

	expectedState := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Step One"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Success"),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.MaxParallelism": types.StringValue("2"),
			"Octopus.Action.TargetRoles":    types.StringValue("agent-1,agent-2"),
		}),
		Type:               types.StringValue("Octopus.Script"),
		Slug:               types.StringValue("step-one"),
		IsRequired:         types.BoolValue(true),
		IsDisabled:         types.BoolValue(false),
		Notes:              types.StringValue(`Some notes`),
		WorkerPoolID:       types.StringValue("WorkerPools-1"),
		WorkerPoolVariable: types.StringValue("#{Environment.WorkerPools.Default}"),
		TenantTags: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("tag-1"),
			types.StringValue("tag-2"),
		}),
		Environments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-1"),
			types.StringValue("Environments-2"),
		}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-13"),
		}),
		Channels: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Channels-1"),
		}),
		Container: &schemas.ProcessStepActionContainerModel{
			FeedID: types.StringValue("Feeds-1"),
			Image:  types.StringValue("docker.io/library/dummy:latest"),
		},
		GitDependencies: types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), map[string]attr.Value{
			gitDependency.Name: types.ObjectValueMust(
				schemas.ProcessStepGitDependencyAttributeTypes(),
				map[string]attr.Value{
					"repository_uri":      types.StringValue(gitDependency.RepositoryUri),
					"default_branch":      types.StringValue(gitDependency.DefaultBranch),
					"git_credential_type": types.StringValue(gitDependency.GitCredentialType),
					"git_credential_id":   types.StringValue(gitDependency.GitCredentialId),
					"file_path_filters": types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("directory-b"),
					}),
				},
			),
		}),
		PrimaryPackage: expectedPrimaryPackage,
		Packages: types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), map[string]attr.Value{
			additionalPackage.Name: types.ObjectValueMust(
				schemas.ProcessStepPackageReferenceAttributeTypes(),
				map[string]attr.Value{
					"id":                   types.StringValue(additionalPackage.ID),
					"package_id":           types.StringValue(additionalPackage.PackageID),
					"feed_id":              types.StringValue(additionalPackage.FeedID),
					"acquisition_location": types.StringValue(additionalPackage.AcquisitionLocation),
					"properties":           types.MapValueMust(types.StringType, map[string]attr.Value{}),
				},
			),
		}),
		ExecutionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.RunOnServer":       types.StringValue("True"),
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Step 1, Action 1\""),
		}),
	}
	expectedState.ID = types.StringValue(step.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapProcessStepToStateWithAllAttributesForRunbook(t *testing.T) {
	primaryPackage := &packages.PackageReference{
		ID:                  "00000000-0000-0000-0000-000000000101",
		Name:                "",
		PackageID:           "Package-1",
		FeedID:              "Feeds-1",
		AcquisitionLocation: "ExecutionTarget",
		Properties: map[string]string{
			"Octopus.Package.Test": "Dummy",
		},
	}
	additionalPackage := &packages.PackageReference{
		ID:                  "00000000-0000-0000-0000-000000000102",
		Name:                "unique-name",
		PackageID:           "Package-2",
		FeedID:              "feeds-builtin",
		AcquisitionLocation: "Server",
	}
	gitDependency := &gitdependencies.GitDependency{
		Name:              "this-dependency",
		RepositoryUri:     "git://test.repository.co.nz",
		DefaultBranch:     "default",
		GitCredentialType: "NotSpecified",
		FilePathFilters:   []string{"directory-b"},
		GitCredentialId:   "GitCredential-2",
	}

	action := deployments.NewDeploymentAction("Step One", "Octopus.Script")
	action.SetID("00000000-0000-0000-0000-000000000011")
	action.Name = "Step One"
	action.Slug = "step-one"
	action.ActionType = "Octopus.Script"
	action.IsRequired = true
	action.IsDisabled = false
	action.Notes = "Some notes"
	action.WorkerPool = "WorkerPools-1"
	action.WorkerPoolVariable = "#{Environment.WorkerPools.Default}"
	action.TenantTags = []string{"tag-1", "tag-2"}
	action.Environments = []string{"Environments-1", "Environments-2"}
	action.ExcludedEnvironments = []string{"Environments-13"}
	action.Channels = []string{"Channels-1"}
	action.Container = &deployments.DeploymentActionContainer{
		FeedID: "Feeds-1",
		Image:  "docker.io/library/dummy:latest",
	}
	action.GitDependencies = []*gitdependencies.GitDependency{gitDependency}
	action.Packages = []*packages.PackageReference{primaryPackage, additionalPackage}
	action.Properties = map[string]core.PropertyValue{
		"Octopus.Action.RunOnServer":                core.NewPropertyValue("True", false),
		"Octopus.Action.Script.ScriptBody":          core.NewPropertyValue("Write-Host \"Step 1, Action 1\"", false),
		"Octopus.Action.Package.FeedId":             core.NewPropertyValue("should_be-ignored-by-primary-package", false),
		"Octopus.Action.Package.PackageId":          core.NewPropertyValue("should_be-ignored-by-primary-package", false),
		"Octopus.Action.Package.DownloadOnTentacle": core.NewPropertyValue("should_be-ignored-by-primary-package", false),
	}

	step := deployments.NewDeploymentStep("Step One")
	step.SetID("00000000-0000-0000-0000-000000000001")
	step.StartTrigger = "StartAfterPrevious"
	step.PackageRequirement = "LetOctopusDecide"
	step.Condition = "Success"
	step.Properties = map[string]core.PropertyValue{
		"Octopus.Action.MaxParallelism": core.NewPropertyValue("2", false),
		"Octopus.Action.TargetRoles":    core.NewPropertyValue("agent-1,agent-2", false),
	}
	step.Actions = []*deployments.DeploymentAction{action}

	process := &runbookprocess.RunbookProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		RunbookID: "Runbooks-2",
		Steps:     []*deployments.DeploymentStep{step},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
	}
	diags := mapProcessStepToState(runbookProcessWrapper{process}, step, &state)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedPrimaryPackage := &schemas.ProcessStepPackageReferenceResourceModel{
		PackageID:           types.StringValue(primaryPackage.PackageID),
		FeedID:              types.StringValue(primaryPackage.FeedID),
		AcquisitionLocation: types.StringValue(primaryPackage.AcquisitionLocation),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Package.Test": types.StringValue("Dummy"),
		}),
	}
	expectedPrimaryPackage.ID = types.StringValue(primaryPackage.ID)

	expectedState := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Step One"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Success"),
		Properties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.MaxParallelism": types.StringValue("2"),
			"Octopus.Action.TargetRoles":    types.StringValue("agent-1,agent-2"),
		}),
		Type:               types.StringValue("Octopus.Script"),
		Slug:               types.StringValue("step-one"),
		IsRequired:         types.BoolValue(true),
		IsDisabled:         types.BoolValue(false),
		Notes:              types.StringValue(`Some notes`),
		WorkerPoolID:       types.StringValue("WorkerPools-1"),
		WorkerPoolVariable: types.StringValue("#{Environment.WorkerPools.Default}"),
		TenantTags: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("tag-1"),
			types.StringValue("tag-2"),
		}),
		Environments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-1"),
			types.StringValue("Environments-2"),
		}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-13"),
		}),
		Channels: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Channels-1"),
		}),
		Container: &schemas.ProcessStepActionContainerModel{
			FeedID: types.StringValue("Feeds-1"),
			Image:  types.StringValue("docker.io/library/dummy:latest"),
		},
		GitDependencies: types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), map[string]attr.Value{
			gitDependency.Name: types.ObjectValueMust(
				schemas.ProcessStepGitDependencyAttributeTypes(),
				map[string]attr.Value{
					"repository_uri":      types.StringValue(gitDependency.RepositoryUri),
					"default_branch":      types.StringValue(gitDependency.DefaultBranch),
					"git_credential_type": types.StringValue(gitDependency.GitCredentialType),
					"git_credential_id":   types.StringValue(gitDependency.GitCredentialId),
					"file_path_filters": types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("directory-b"),
					}),
				},
			),
		}),
		PrimaryPackage: expectedPrimaryPackage,
		Packages: types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), map[string]attr.Value{
			additionalPackage.Name: types.ObjectValueMust(
				schemas.ProcessStepPackageReferenceAttributeTypes(),
				map[string]attr.Value{
					"id":                   types.StringValue(additionalPackage.ID),
					"package_id":           types.StringValue(additionalPackage.PackageID),
					"feed_id":              types.StringValue(additionalPackage.FeedID),
					"acquisition_location": types.StringValue(additionalPackage.AcquisitionLocation),
					"properties":           types.MapValueMust(types.StringType, map[string]attr.Value{}),
				},
			),
		}),
		ExecutionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.RunOnServer":       types.StringValue("True"),
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Step 1, Action 1\""),
		}),
	}
	expectedState.ID = types.StringValue(step.ID)

	assert.Equal(t, expectedState, state)
}
