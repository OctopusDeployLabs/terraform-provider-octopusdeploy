package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantModel struct {
	ClonedFromTenantId types.String `tfsdk:"cloned_from_tenant_id"`
	Description        types.String `tfsdk:"description"`
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	SpaceID            types.String `tfsdk:"space_id"`
	TenantTags         types.List   `tfsdk:"tenant_tags"`
}

type TenantsModel struct {
	ClonedFromTenantId types.String `tfsdk:"cloned_from_tenant_id"`
	ID                 types.String `tfsdk:"id"`
	IDs                types.List   `tfsdk:"ids"`
	IsClone            types.Bool   `tfsdk:"is_clone"`
	Name               types.String `tfsdk:"name"`
	PartialName        types.String `tfsdk:"partial_name"`
	ProjectId          types.String `tfsdk:"project_id"`
	Skip               types.Int64  `tfsdk:"skip"`
	Tags               types.List   `tfsdk:"tags"`
	SpaceID            types.String `tfsdk:"space_id"`
	Tenants            types.List   `tfsdk:"tenants"`
	Take               types.Int64  `tfsdk:"take"`
}

func TenantObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"cloned_from_tenant_id": types.StringType,
		"description":           types.StringType,
		"id":                    types.StringType,
		"name":                  types.StringType,
		"space_id":              types.StringType,
		"tenant_tags":           types.ListType{ElemType: types.StringType},
	}
}

func FlattenTenant(tenant *tenants.Tenant) attr.Value {
	tenantTags := make([]attr.Value, len(tenant.TenantTags))
	for i, value := range tenant.TenantTags {
		tenantTags[i] = types.StringValue(value)
	}
	var tenantTagsList, _ = types.ListValue(types.StringType, tenantTags)

	return types.ObjectValueMust(TenantObjectType(), map[string]attr.Value{
		"cloned_from_tenant_id": types.StringValue(tenant.ClonedFromTenantID),
		"description":           types.StringValue(tenant.Description),
		"id":                    types.StringValue(tenant.GetID()),
		"name":                  types.StringValue(tenant.Name),
		"space_id":              types.StringValue(tenant.SpaceID),
		"tenant_tags":           tenantTagsList,
	})
}

func GetTenantsDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"cloned_from_tenant_id": datasourceSchema.StringAttribute{
			Description: "A filter to search for a cloned tenant by its ID.",
			Optional:    true,
		},
		"id":  util.GetIdDatasourceSchema(),
		"ids": util.GetQueryIDsDatasourceSchema(),
		"is_clone": datasourceSchema.BoolAttribute{
			Description: "A filter to search for cloned resources.",
			Optional:    true,
		},
		"name": datasourceSchema.StringAttribute{
			Description: "A filter to search by name.",
			Optional:    true,
		},
		"partial_name": util.GetQueryPartialNameDatasourceSchema(),
		"project_id": datasourceSchema.StringAttribute{
			Description: "A filter to search by a project ID.",
			Optional:    true,
		},
		"skip":     util.GetQuerySkipDatasourceSchema(),
		"tags":     util.GetQueryDatasourceTags(),
		"space_id": util.GetSpaceIdDatasourceSchema("tenants"),
		"take":     util.GetQueryTakeDatasourceSchema(),
	}
}

func GetTenantDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"cloned_from_tenant_id": datasourceSchema.StringAttribute{
			Description: "The ID of the tenant from which this tenant was cloned.",
			Optional:    true,
		},
		"description": util.GetDescriptionDatasourceSchema("tenants"),
		"id":          util.GetIdDatasourceSchema(),
		"name":        util.GetNameDatasourceSchema(false),
		"space_id":    util.GetSpaceIdDatasourceSchema("tenant"),
		"tenant_tags": datasourceSchema.ListAttribute{
			Computed:    true,
			Description: "A list of tenant tags associated with this resource.",
			ElementType: types.StringType,
			Optional:    true,
		},
	}
}

func GetTenantResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"cloned_from_tenant_id": resourceSchema.StringAttribute{
			Description: "The ID of the tenant from which this tenant was cloned.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"description": util.GetDescriptionResourceSchema("tenant"),
		"id":          util.GetIdResourceSchema(),
		"name":        util.GetNameResourceSchema(true),
		"space_id":    util.GetSpaceIdResourceSchema("tenant"),
		"tenant_tags": resourceSchema.ListAttribute{
			Computed:    true,
			Description: "A list of tenant tags associated with this resource.",
			ElementType: types.StringType,
			Optional:    true,
		},
	}
}
