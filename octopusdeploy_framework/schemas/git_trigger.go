package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GitTriggerSchema struct{}

var _ EntitySchema = GitTriggerSchema{}

func (d GitTriggerSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages Git triggers in Octopus Deploy",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema("Git trigger."),
			"space_id":    GetSpaceIdResourceSchema("Git trigger"),
			"project_id":  GetRequiredStringResourceSchema("The ID of the project to attach the trigger."),
			"channel_id":  GetRequiredStringResourceSchema("The ID of the channel in which the release will be created if the action type is CreateRelease."),
			"sources":     GetSourcesAttributeSchema(),
			"is_disabled": GetOptionalBooleanResourceAttribute("Disables the trigger from being run when set.", false),
		},
	}
}

func GetSourcesAttributeSchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Required: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"deployment_action_slug": util.DataSourceString().Required().Description("The deployment action slug.").Build(),
				"git_dependency_name":    util.DataSourceString().Required().Description("The git dependency name.").Build(),
				"include_file_paths": resourceSchema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
					Description: "The file paths to include.",
				},
				"exclude_file_paths": resourceSchema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
					Description: "The file paths to exclude.",
				},
			},
		},
	}
}

func (d GitTriggerSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type GitTriggerResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SpaceId     types.String `tfsdk:"space_id"`
	ProjectId   types.String `tfsdk:"project_id"`
	ChannelId   types.String `tfsdk:"channel_id"`
	Sources     types.List   `tfsdk:"sources"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`

	ResourceModel
}
