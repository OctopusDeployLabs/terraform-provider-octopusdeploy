package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentFreezeSchema struct{}

func (d DeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":    GetIdResourceSchema(),
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("The start time of the freeze, must be RFC3339 format", true),
			"end":   GetDateTimeResourceSchema("The end time of the freeze, must be RFC3339 format", true),
		},
	}
}

func (d DeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about deployment freezes",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           GetIdDatasourceSchema(true),
			"ids":          GetQueryIDsDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"project_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by a list of project IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"environment_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by a list of environment IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"include_complete": GetBooleanDatasourceAttribute("Include deployment freezes that completed, default is true", true),
			"status": datasourceSchema.StringAttribute{
				Description: "Filter by the status of the deployment freeze, value values are Expired, Active, Scheduled (case-insensitive)",
				Optional:    true,
			},
			"deployment_freezes": datasourceSchema.ListNestedAttribute{
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":   GetIdDatasourceSchema(true),
						"name": GetReadonlyNameDatasourceSchema(),
						"start": datasourceSchema.StringAttribute{
							Description: "The start time of the freeze",
							Optional:    false,
							Computed:    true,
						},
						"end": datasourceSchema.StringAttribute{
							Description: "The end time of the freeze",
							Optional:    false,
							Computed:    true,
						},
						"project_environment_scope": datasourceSchema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Description: "The project environment scope of the deployment freeze",
							Optional:    false,
							Computed:    true,
						},
					},
				},
				Optional: false,
				Computed: true,
			},
		},
	}

}

var _ EntitySchema = &DeploymentFreezeSchema{}
