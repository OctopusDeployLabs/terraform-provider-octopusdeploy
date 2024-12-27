package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BuiltInTriggerSchema struct{}

var _ EntitySchema = BuiltInTriggerSchema{}

func (r BuiltInTriggerSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"project_id": util.ResourceString().
				Description("The ID of the project the trigger will be attached to.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Required().
				Build(),
			"space_id": util.ResourceString().
				Description("Space ID of the associated project.").
				Optional().
				Computed().
				Build(),
			"channel_id": util.ResourceString().
				Description("The ID of the channel in which triggered release will be created.").
				Required().
				Build(),
			"release_creation_package_step_id": util.ResourceString().
				Description("The package step ID trigger will be listening.").
				Optional().
				Computed().
				Build(),
			"release_creation_package": resourceSchema.SingleNestedAttribute{
				Required:    true,
				Description: "Combination of deployment action and package references.",
				Attributes: map[string]resourceSchema.Attribute{
					"deployment_action": util.ResourceString().
						Description("Deployment action.").
						Optional().
						Build(),
					"package_reference": util.ResourceString().
						Description("Package reference.").
						Optional().
						Build(),
				},
			},
		},
	}
}

func (r BuiltInTriggerSchema) GetDatasourceSchema() datasourceSchema.Schema {
	// no datasource required, returned as part of project datasource
	return datasourceSchema.Schema{}
}

type BuiltInTriggerResourceModel struct {
	ProjectID                    types.String                `tfsdk:"project_id"`
	SpaceID                      types.String                `tfsdk:"space_id"`
	ChannelID                    types.String                `tfsdk:"channel_id"`
	ReleaseCreationPackageStepID types.String                `tfsdk:"release_creation_package_step_id"`
	ReleaseCreationPackage       ReleaseCreationPackageModel `tfsdk:"release_creation_package"`
}

type ReleaseCreationPackageModel struct {
	DeploymentAction types.String `tfsdk:"deployment_action"`
	PackageReference types.String `tfsdk:"package_reference"`
}
