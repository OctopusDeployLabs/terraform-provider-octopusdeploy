package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessChildStepsOrderSchema struct{}

var _ EntitySchema = ProcessChildStepsOrderSchema{}

const ProcessChildStepsOrderResourceName = "process_child_steps_order"

func (p ProcessChildStepsOrderSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages order of Child Steps in a Runbook or Deployment Process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessChildStepsOrderResourceName),
			"process_id": util.ResourceString().
				Description("Id of the process parent step belongs to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"parent_id": util.ResourceString().
				Description("Id of the process step children belong to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"children": util.ResourceList(types.StringType).
				Description("Child steps in the order of execution").
				Required().
				Validators(listvalidator.UniqueValues()).
				Build(),
		},
	}
}

func (p ProcessChildStepsOrderSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessChildStepsOrderResourceModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	ProcessID types.String `tfsdk:"process_id"`
	ParentID  types.String `tfsdk:"parent_id"`
	Children  types.List   `tfsdk:"children"`

	ResourceModel
}
