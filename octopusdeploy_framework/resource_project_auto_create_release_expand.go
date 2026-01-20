package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func expand(ctx context.Context, data *schemas.ProjectAutoCreateReleaseResourceModel) *projects.ReleaseCreationStrategy {
	if len(data.ReleaseCreationPackage) == 0 {
		return &projects.ReleaseCreationStrategy{}
	}

	pkg := data.ReleaseCreationPackage[0]
	strategy := &projects.ReleaseCreationStrategy{
		ChannelID:                    data.ChannelID.ValueString(),
		ReleaseCreationPackageStepID: data.ReleaseCreationPackageStepID.ValueString(),
		ReleaseCreationPackage: &packages.DeploymentActionPackage{
			DeploymentAction: pkg.DeploymentAction.ValueString(),
			PackageReference: pkg.PackageReference.ValueString(),
		},
	}

	return strategy
}

func flatten(ctx context.Context, strategy *projects.ReleaseCreationStrategy, data *schemas.ProjectAutoCreateReleaseResourceModel) {
	if strategy == nil || strategy.ReleaseCreationPackage == nil {
		data.ChannelID = types.StringNull()
		data.ReleaseCreationPackageStepID = types.StringNull()
		data.ReleaseCreationPackage = []schemas.ProjectAutoCreateReleaseCreationPackage{}
		return
	}

	data.ChannelID = types.StringValue(strategy.ChannelID)
	data.ReleaseCreationPackageStepID = types.StringValue(strategy.ReleaseCreationPackageStepID)
	data.ReleaseCreationPackage = []schemas.ProjectAutoCreateReleaseCreationPackage{
		{
			DeploymentAction: types.StringValue(strategy.ReleaseCreationPackage.DeploymentAction),
			PackageReference: types.StringValue(strategy.ReleaseCreationPackage.PackageReference),
		},
	}
}
