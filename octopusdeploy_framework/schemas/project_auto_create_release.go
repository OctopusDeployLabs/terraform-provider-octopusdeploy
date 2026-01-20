package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectAutoCreateReleaseSchema struct{}

var _ EntitySchema = ProjectAutoCreateReleaseSchema{}

func (r ProjectAutoCreateReleaseSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages automatic release creation for a project. When enabled, new releases will be automatically created when new versions of the specified packages are available. The referenced packages must use built-in package repository feeds.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": util.ResourceString().
				Description("The unique identifier for this resource.").
				Computed().
				Build(),
			"project_id": util.ResourceString().
				Description("The ID of the project for which to enable automatic release creation.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Required().
				Build(),
			"space_id": util.ResourceString().
				Description("The space ID where the project is located. If not specified, the default space will be used.").
				Optional().
				Computed().
				Build(),
			"channel_id": util.ResourceString().
				Description("The ID of the channel in which triggered releases will be created.").
				Required().
				Build(),
			"release_creation_package_step_id": util.ResourceString().
				Description("The ID of the deployment step containing the package for release creation.").
				Optional().
				Computed().
				Build(),
		},
		Blocks: map[string]resourceSchema.Block{
			"release_creation_package": resourceSchema.ListNestedBlock{
				Description: "Configuration for the package that will trigger automatic release creation. The referenced package must use a built-in package repository feed.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"deployment_action": util.ResourceString().
							Description("The name of the deployment action that contains the package reference.").
							Required().
							Build(),
						"package_reference": util.ResourceString().
							Description("The name of the package reference within the deployment action.").
							Required().
							Build(),
					},
				},
			},
		},
	}
}

func (r ProjectAutoCreateReleaseSchema) GetDatasourceSchema() datasourceSchema.Schema {
	// No datasource required for this resource
	return datasourceSchema.Schema{}
}

type ProjectAutoCreateReleaseResourceModel struct {
	ID                           types.String                              `tfsdk:"id"`
	ProjectID                    types.String                              `tfsdk:"project_id"`
	SpaceID                      types.String                              `tfsdk:"space_id"`
	ChannelID                    types.String                              `tfsdk:"channel_id"`
	ReleaseCreationPackageStepID types.String                              `tfsdk:"release_creation_package_step_id"`
	ReleaseCreationPackage       []ProjectAutoCreateReleaseCreationPackage `tfsdk:"release_creation_package"`
}

type ProjectAutoCreateReleaseCreationPackage struct {
	DeploymentAction types.String `tfsdk:"deployment_action"`
	PackageReference types.String `tfsdk:"package_reference"`
}
