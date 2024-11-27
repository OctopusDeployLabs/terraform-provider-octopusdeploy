package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentFreezeProjectSchema struct{}

type DeploymentFreezeProjectResourceModel struct {
	DeploymentFreezeID types.String `tfsdk:"deploymentfreeze_id"`
	ProjectID          types.String `tfsdk:"project_id"`
	EnvironmentIDs     types.List   `tfsdk:"environment_ids"`
	ResourceModel
}

func (d DeploymentFreezeProjectSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"deploymentfreeze_id": resourceSchema.StringAttribute{
				Description: "The deployment freeze ID associated with this freeze scope.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": resourceSchema.StringAttribute{
				Description: "The project ID associated with this freeze scope.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_ids": resourceSchema.ListAttribute{
				Description: "The environment IDs associated with this project deployment freeze scope.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (d DeploymentFreezeProjectSchema) GetDatasourceSchema() datasourceSchema.Schema {
	//TODO implement me
	panic("implement me")
}

var _ EntitySchema = DeploymentFreezeProjectSchema{}
