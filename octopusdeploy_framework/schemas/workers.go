package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkersSchema struct{}

var _ EntitySchema = WorkersSchema{}

func (f WorkersSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

func (f WorkersSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing workers.",
		Attributes: map[string]datasourceSchema.Attribute{
			"ids":          GetQueryIDsDatasourceSchema(),
			"name":         GetNameDatasourceSchema(false),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"space_id":     GetSpaceIdDatasourceSchema("workers", false),
			"communication_styles": datasourceSchema.ListAttribute{
				Description: "A filter to search by communication styles",
				ElementType: types.StringType,
				Optional:    true,
			},
			"health_statuses": datasourceSchema.ListAttribute{
				Description: "A filter to search by health statuses",
				ElementType: types.StringType,
				Optional:    true,
			},
			"worker_pool_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by worker pools",
				ElementType: types.StringType,
				Optional:    true,
			},
			"is_disabled": GetBooleanDatasourceAttribute("", true),
			"thumbprint": datasourceSchema.StringAttribute{
				Description: "A filter search by worker's thumbprint",
				Optional:    true,
			},

			// response
			"id": GetIdDatasourceSchema(true),
			"workers": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":       GetIdDatasourceSchema(true),
						"space_id": GetSpaceIdDatasourceSchema("feeds", true),
						"name":     GetReadonlyNameDatasourceSchema(),
						"is_disabled": datasourceSchema.BoolAttribute{
							Computed: true,
						},
						"machine_policy_id": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"worker_pool_ids": datasourceSchema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"proxy_id": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"uri": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"thumbprint": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"account_id": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"host": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"port": datasourceSchema.Int64Attribute{
							Computed: true,
						},
						"fingerprint": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"dotnet_platform": datasourceSchema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type WorkersDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	IDs                types.List   `tfsdk:"ids"`
	Name               types.String `tfsdk:"name"`
	PartialName        types.String `tfsdk:"partial_name"`
	Skip               types.Int64  `tfsdk:"skip"`
	Take               types.Int64  `tfsdk:"take"`
	SpaceID            types.String `tfsdk:"space_id"`
	CommunicationStyle types.List   `tfsdk:"health_statuses"`
	HealthStatuses     types.List   `tfsdk:"communication_styles"`
	WorkerPoolIDs      types.List   `tfsdk:"worker_pool_ids"`
	IsDisabled         types.Bool   `tfsdk:"is_disabled"`
	Thumbprint         types.String `tfsdk:"thumbprint"`
	Workers            types.List   `tfsdk:"workers"`
}

func FlattenWorker(worker *machines.Worker) attr.Value {
	return types.ObjectValueMust(WorkerObjectType(), map[string]attr.Value{
		"id":                  types.StringValue(worker.GetID()),
		"name":                types.StringValue(worker.Name),
		"is_disabled":         types.BoolValue(worker.IsDisabled),
		"communication_style": types.StringValue(worker.Endpoint.GetCommunicationStyle()),
		"health_status":       types.StringValue(worker.HealthStatus),
		"machine_policy_id":   types.StringValue(worker.MachinePolicyID),
		"worker_pool_ids":     types.ListValueMust(types.StringType, util.ToValueSlice(worker.WorkerPoolIDs)),
	})
}

func WorkerObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.StringType,
		"name":                types.StringType,
		"is_disabled":         types.BoolType,
		"communication_style": types.StringType,
		"health_status":       types.StringType,
		"machine_policy_id":   types.StringType,
		"worker_pool_ids":     types.ListType{ElemType: types.StringType},
		//"proxy_id":            types.StringType, // Endpoint specific values
		//"uri":                 types.StringType,
		//"thumbprint":          types.StringType,
		//"account_id":      types.StringType,
		//"host":            types.StringType,
		//"port":            types.Int64Type,
		//"fingerprint":     types.StringType,
		//"dotnet_platform": types.StringType,
	}
}
