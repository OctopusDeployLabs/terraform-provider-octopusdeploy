package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantModel struct {
	ClonedFromTenantId types.String `tfsdk:"cloned_from_tenant_id"`
	Description        types.String `tfsdk:"description"`
	IsDisabled         types.Bool   `tfsdk:"is_disabled"`
	Name               types.String `tfsdk:"name"`
	SpaceID            types.String `tfsdk:"space_id"`
	TenantTags         types.Set    `tfsdk:"tenant_tags"`

	ResourceModel
}

type TenantsModel struct {
	ClonedFromTenantId types.String `tfsdk:"cloned_from_tenant_id"`
	ID                 types.String `tfsdk:"id"`
	IDs                types.List   `tfsdk:"ids"`
	IsClone            types.Bool   `tfsdk:"is_clone"`
	IsDisabled         types.Bool   `tfsdk:"is_disabled"`
	Name               types.String `tfsdk:"name"`
	PartialName        types.String `tfsdk:"partial_name"`
	ProjectId          types.String `tfsdk:"project_id"`
	Skip               types.Int64  `tfsdk:"skip"`
	Tags               types.List   `tfsdk:"tags"`
	SpaceID            types.String `tfsdk:"space_id"`
	Tenants            types.List   `tfsdk:"tenants"`
	Take               types.Int64  `tfsdk:"take"`
}

type TenantSchema struct{}

var _ EntitySchema = TenantSchema{}

func TenantObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"cloned_from_tenant_id": types.StringType,
		"description":           types.StringType,
		"id":                    types.StringType,
		"is_disabled":           types.BoolType,
		"name":                  types.StringType,
		"space_id":              types.StringType,
		"tenant_tags":           types.SetType{ElemType: types.StringType},
	}
}

func FlattenTenant(tenant *tenants.Tenant) attr.Value {
	tenantTags := make([]attr.Value, len(tenant.TenantTags))
	for i, value := range tenant.TenantTags {
		tenantTags[i] = types.StringValue(value)
	}
	var tenantTagsSet, _ = types.SetValue(types.StringType, tenantTags)

	return types.ObjectValueMust(TenantObjectType(), map[string]attr.Value{
		"cloned_from_tenant_id": types.StringValue(tenant.ClonedFromTenantID),
		"description":           types.StringValue(tenant.Description),
		"id":                    types.StringValue(tenant.GetID()),
		"is_disabled":           types.BoolValue(tenant.IsDisabled),
		"name":                  types.StringValue(tenant.Name),
		"space_id":              types.StringValue(tenant.SpaceID),
		"tenant_tags":           tenantTagsSet,
	})
}

func (t TenantSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing tenants.",
		Attributes: map[string]datasourceSchema.Attribute{
			"cloned_from_tenant_id": datasourceSchema.StringAttribute{
				Description: "A filter to search for a cloned tenant by its ID.",
				Optional:    true,
			},
			"id":  GetIdDatasourceSchema(true),
			"ids": GetQueryIDsDatasourceSchema(),
			"is_clone": datasourceSchema.BoolAttribute{
				Description: "A filter to search for cloned resources.",
				Optional:    true,
			},
			"is_disabled": datasourceSchema.BoolAttribute{
				Description: "A filter to search by the disabled status of a resource.",
				Optional:    true,
			},
			"name": datasourceSchema.StringAttribute{
				Description: "A filter to search by name.",
				Optional:    true,
			},
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"project_id": datasourceSchema.StringAttribute{
				Description: "A filter to search by a project ID.",
				Optional:    true,
			},
			"skip":     GetQuerySkipDatasourceSchema(),
			"tags":     GetQueryDatasourceTags(),
			"space_id": GetSpaceIdDatasourceSchema("tenants", false),
			"take":     GetQueryTakeDatasourceSchema(),
			"tenants": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"cloned_from_tenant_id": datasourceSchema.StringAttribute{
							Description: "The ID of the tenant from which this tenant was cloned.",
							Computed:    true,
						},
						"description": GetDescriptionDatasourceSchema("tenants"),
						"id":          GetIdDatasourceSchema(true),
						"is_disabled": datasourceSchema.BoolAttribute{
							Description: "The disabled status of this tenant.",
							Computed:    true,
						},
						"name":     GetReadonlyNameDatasourceSchema(),
						"space_id": GetSpaceIdDatasourceSchema("tenant", true),
						"tenant_tags": datasourceSchema.SetAttribute{
							Computed:    true,
							Description: "A list of tenant tags associated with this resource.",
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (t TenantSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages tenants in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"cloned_from_tenant_id": resourceSchema.StringAttribute{
				Description: "The ID of the tenant from which this tenant was cloned.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": GetDescriptionResourceSchema("tenant"),
			"id":          GetIdResourceSchema(),
			"is_disabled": datasourceSchema.BoolAttribute{
				Description: "The disabled status of this tenant.",
				Optional:    true,
				Computed:    true,
			},
			"name":     GetNameResourceSchema(true),
			"space_id": GetSpaceIdResourceSchema("tenant"),
			"tenant_tags": resourceSchema.SetAttribute{
				Description: "A list of tenant tags associated with this resource.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
