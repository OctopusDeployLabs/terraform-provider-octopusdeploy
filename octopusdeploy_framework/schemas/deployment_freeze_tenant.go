package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentFreezeTenantSchema struct{}

type DeploymentFreezeTenantResourceModel struct {
	DeploymentFreezeID types.String `tfsdk:"deploymentfreeze_id"`
	TenantID           types.String `tfsdk:"tenant_id"`
	ProjectID          types.String `tfsdk:"project_id"`
	EnvironmentID      types.String `tfsdk:"environment_id"`
	ResourceModel
}

func (d DeploymentFreezeTenantSchema) GetResourceSchema() resourceSchema.Schema {
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
			"tenant_id": resourceSchema.StringAttribute{
				Description: "The tenant ID associated with this freeze scope.",
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
			"environment_id": resourceSchema.StringAttribute{
				Description: "The environment ID associated with this freeze scope.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (d DeploymentFreezeTenantSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about deployment freeze tenant scopes",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":  GetIdDatasourceSchema(true),
			"ids": GetQueryIDsDatasourceSchema(),
			"deploymentfreeze_id": datasourceSchema.StringAttribute{
				Description: "The deployment freeze ID associated with this freeze scope",
				Required:    true,
			},
			"tenant_id": datasourceSchema.StringAttribute{
				Description: "The tenant ID associated with this freeze scope",
				Required:    true,
			},
			"project_id": datasourceSchema.StringAttribute{
				Description: "The project ID associated with this freeze scope",
				Required:    true,
			},
			"environment_id": datasourceSchema.StringAttribute{
				Description: "The environment ID associated with this freeze scope",
				Required:    true,
			},
		},
	}
}

var _ EntitySchema = DeploymentFreezeTenantSchema{}
