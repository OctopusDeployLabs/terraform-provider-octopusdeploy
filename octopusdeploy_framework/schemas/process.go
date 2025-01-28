package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProcessSchema struct{}

var _ EntitySchema = ProcessSchema{}

const ProcessResourceName = "process"

func (p ProcessSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages execution processes in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(ProcessResourceName),
			"owner_id": util.ResourceString().Required().Description("Id of the resource this process belongs to.").Build(),
		},
	}
}

func (p ProcessSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type ProcessResourceModel struct {
	SpaceID types.String `tfsdk:"space_id"`
	OwnerID types.String `tfsdk:"owner_id"`

	ResourceModel
}
