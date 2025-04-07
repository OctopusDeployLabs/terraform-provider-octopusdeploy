package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessSchema struct{}

var _ EntitySchema = ProcessSchema{}

const ProcessResourceName = "process"

func (p ProcessSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages Runbook and Deployment Processes in Octopus Deploy. It's used in collaboration with `octopusdeploy_process_step` and `octopusdeploy_process_step_order`.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessResourceName),
			"project_id": util.ResourceString().
				Required().
				Description("Id of the project this process belongs to.").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"runbook_id": util.ResourceString().
				Optional().
				Description("Id of the runbook this process belongs to. When not set this resource represents deployment process of the project").
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
		},
	}
}

func (p ProcessSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessResourceModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	ProjectID types.String `tfsdk:"project_id"`
	RunbookID types.String `tfsdk:"runbook_id"`

	ResourceModel
}
