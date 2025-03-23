package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccMapBuiltInTriggerFromState(t *testing.T) {
	project := &projects.Project{
		SpaceID:                 "Spaces-1",
		Name:                    "Test With Trigger",
		AutoCreateRelease:       false,
		ReleaseCreationStrategy: nil,
	}
	project.SetID("Projects-21")

	state := schemas.BuiltInTriggerResourceModel{
		SpaceID:                      types.StringValue(project.SpaceID),
		ProjectID:                    types.StringValue(project.ID),
		ChannelID:                    types.StringValue("Channels-42"),
		ReleaseCreationPackageStepID: types.StringValue("10000000-0000-0000-0000-000000000001"),
		ReleaseCreationPackage: schemas.ReleaseCreationPackageModel{
			DeploymentAction: types.StringValue("Name of the Action"),
			PackageReference: types.StringValue("my-package"),
		},
	}

	mapBuiltInTriggerFromState(&state, project)

	expected := &projects.Project{
		SpaceID:           project.SpaceID,
		Name:              project.Name,
		AutoCreateRelease: true,
		ReleaseCreationStrategy: &projects.ReleaseCreationStrategy{
			ChannelID:                    "Channels-42",
			ReleaseCreationPackageStepID: "10000000-0000-0000-0000-000000000001",
			ReleaseCreationPackage: &packages.DeploymentActionPackage{
				DeploymentAction: "Name of the Action",
				PackageReference: "my-package",
			},
		},
	}
	expected.SetID(project.ID)

	assert.Equal(t, expected, project)
}

func TestAccMapBuiltInTriggerToState(t *testing.T) {
	project := &projects.Project{
		SpaceID:           "Spaces-2",
		Name:              "Map to state",
		AutoCreateRelease: true,
		ReleaseCreationStrategy: &projects.ReleaseCreationStrategy{
			ChannelID:                    "Channels-82",
			ReleaseCreationPackageStepID: "10000000-0000-0000-0000-000000000002",
			ReleaseCreationPackage: &packages.DeploymentActionPackage{
				DeploymentAction: "Map",
				PackageReference: "map-package",
			},
		},
	}
	project.SetID("Projects-31")

	state := &schemas.BuiltInTriggerResourceModel{
		SpaceID:                      types.StringValue(project.SpaceID),
		ProjectID:                    types.StringValue(project.ID),
		ChannelID:                    types.StringValue("Channels-122"),
		ReleaseCreationPackageStepID: types.StringNull(),
		ReleaseCreationPackage: schemas.ReleaseCreationPackageModel{
			DeploymentAction: types.StringNull(),
			PackageReference: types.StringNull(),
		},
	}

	expected := &schemas.BuiltInTriggerResourceModel{
		SpaceID:                      types.StringValue(project.SpaceID),
		ProjectID:                    types.StringValue(project.ID),
		ChannelID:                    types.StringValue("Channels-82"),
		ReleaseCreationPackageStepID: types.StringValue("10000000-0000-0000-0000-000000000002"),
		ReleaseCreationPackage: schemas.ReleaseCreationPackageModel{
			DeploymentAction: types.StringValue("Map"),
			PackageReference: types.StringValue("map-package"),
		},
	}

	diagnostics := mapBuiltInTriggerToState(project, state)

	assert.Equal(t, 0, len(diagnostics), "Expected mapping diagnostics to be empty")
	assert.Equal(t, expected, state)
}

func TestAccMapBuiltInTriggerToStateWithoutStrategy(t *testing.T) {
	project := &projects.Project{
		SpaceID:                 "Spaces-2",
		Name:                    "Map to state",
		AutoCreateRelease:       false,
		ReleaseCreationStrategy: nil,
	}
	project.SetID("Projects-31")

	state := &schemas.BuiltInTriggerResourceModel{
		SpaceID:                      types.StringValue(project.SpaceID),
		ProjectID:                    types.StringValue(project.ID),
		ChannelID:                    types.StringValue("Channels-122"),
		ReleaseCreationPackageStepID: types.StringNull(),
		ReleaseCreationPackage: schemas.ReleaseCreationPackageModel{
			DeploymentAction: types.StringNull(),
			PackageReference: types.StringNull(),
		},
	}

	diagnostics := mapBuiltInTriggerToState(project, state)

	expected := &schemas.BuiltInTriggerResourceModel{
		SpaceID:                      types.StringValue(project.SpaceID),
		ProjectID:                    types.StringValue(project.ID),
		ChannelID:                    types.StringValue("Channels-82"),
		ReleaseCreationPackageStepID: types.StringValue("10000000-0000-0000-0000-000000000002"),
		ReleaseCreationPackage: schemas.ReleaseCreationPackageModel{
			DeploymentAction: types.StringValue("Map"),
			PackageReference: types.StringValue("map-package"),
		},
	}

	assert.Equal(t, 0, len(diagnostics), "Expected mapping diagnostics to be empty")
	assert.Equal(t, expected, state)
}
