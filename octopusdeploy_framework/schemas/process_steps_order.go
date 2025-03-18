package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessStepsOrderSchema struct{}

var _ EntitySchema = ProcessStepsOrderSchema{}

const ProcessStepsOrderResourceName = "process_steps_order"

func (p ProcessStepsOrderSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages order of steps in the process.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":         GetIdResourceSchema(),
			"space_id":   GetSpaceIdResourceSchema(ProcessStepsOrderResourceName),
			"process_id": util.ResourceString().Required().Description("Id of the process steps belongs to.").Build(),
			"steps": util.ResourceList(types.StringType).
				Description("Steps in the order of execution").
				Required().
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
