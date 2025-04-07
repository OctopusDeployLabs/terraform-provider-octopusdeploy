package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessStepsOrderSchema struct{}

var _ EntitySchema = ProcessStepsOrderSchema{}

const ProcessStepsOrderResourceName = "process_steps_order"

func (p ProcessStepsOrderSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages order of Steps of a Runbook or Deployment Process in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessStepsOrderResourceName),
			"process_id": util.ResourceString().
				Description("Id of the process steps belongs to.").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"steps": util.ResourceList(types.StringType).
				Description("Steps in the order of execution").
				Required().
				Validators(listvalidator.UniqueValues()).
				Build(),
		},
	}
}

func (p ProcessStepsOrderSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessStepsOrderResourceModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	ProcessID types.String `tfsdk:"process_id"`
	Steps     types.List   `tfsdk:"steps"`

	ResourceModel
}
