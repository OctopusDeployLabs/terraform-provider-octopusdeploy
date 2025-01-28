package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessChildStepSchema struct{}

var _ EntitySchema = ProcessChildStepSchema{}

const ProcessChildStepResourceName = "process_child_step"

func (p ProcessChildStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a child step in execution process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":         GetIdResourceSchema(),
			"space_id":   GetSpaceIdResourceSchema(ProcessChildStepResourceName),
			"process_id": util.ResourceString().Required().Description("Id of the process this step belongs to.").Build(),
			"parent_id":  util.ResourceString().Required().Description("Id of the process step this step belongs to.").Build(),
			"name":       GetNameResourceSchema(true),
			"action_type": util.ResourceString().
				Description("Type of the step action").
				Required().
				Build(),
			"run_on_server": util.ResourceBool().
				Description("Whether this step runs on a worker or on the target").
				Default(true).
				Optional().
				Computed().
				Build(),
			"script_syntax": util.ResourceString().
				Description("Type of the syntax of the script").
				Optional().
				Computed().
				Build(),
			"script_body": util.ResourceString().
				Description("Script to execute").
				Optional().
				Build(),
		},
	}
}

func (p ProcessChildStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessChildStepResourceModel struct {
	SpaceID      types.String `tfsdk:"space_id"`
	ProcessID    types.String `tfsdk:"process_id"`
	ParentID     types.String `tfsdk:"parent_id"`
	Name         types.String `tfsdk:"name"`
	ActionType   types.String `tfsdk:"action_type"`
	RunOnServer  types.Bool   `tfsdk:"run_on_server"`
	ScriptSyntax types.String `tfsdk:"script_syntax"`
	ScriptBody   types.String `tfsdk:"script_body"`

	ResourceModel
}
