package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessStepSchema struct{}

var _ EntitySchema = ProcessStepSchema{}

const ProcessStepResourceName = "process_step"

func (p ProcessStepSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages single step of execution process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":         GetIdResourceSchema(),
			"space_id":   GetSpaceIdResourceSchema(ProcessStepResourceName),
			"process_id": util.ResourceString().Required().Description("Id of the process this step belongs to.").Build(),
			"name":       GetNameResourceSchema(true),
			"condition": util.ResourceString().
				Description("When to run the step, one of 'Success', 'Failure', 'Always' or 'Variable'").
				Optional().
				Computed().
				Default("Success").
				Validators(stringvalidator.OneOf("Success", "Failure", "Always", "Variable")).
				Build(),
			"start_trigger": util.ResourceString().
				Description("Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')").
				Optional().
				Computed().
				Default("StartAfterPrevious").
				Validators(stringvalidator.OneOf("StartAfterPrevious", "StartWithPrevious")).
				Build(),
			"target_roles": util.ResourceSet(types.StringType).
				Description("The roles that this step run against, or runs on behalf of").
				Optional().
				Build(),
			"window_size": util.ResourceString().
				Description("The maximum number of targets to deploy to simultaneously").
				Optional().
				Build(),
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

func (p ProcessStepSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessStepResourceModel struct {
	SpaceID      types.String `tfsdk:"space_id"`
	ProcessID    types.String `tfsdk:"process_id"`
	Name         types.String `tfsdk:"name"`
	Condition    types.String `tfsdk:"condition"`
	StartTrigger types.String `tfsdk:"start_trigger"`
	TargetRoles  types.Set    `tfsdk:"target_roles"`
	WindowSize   types.String `tfsdk:"window_size"`
	ActionType   types.String `tfsdk:"action_type"`
	ScriptSyntax types.String `tfsdk:"script_syntax"`
	ScriptBody   types.String `tfsdk:"script_body"`
	RunOnServer  types.Bool   `tfsdk:"run_on_server"`

	ResourceModel
}
