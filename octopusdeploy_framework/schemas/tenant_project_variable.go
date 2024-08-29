package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	TenantProjectVariableResourceDescription = "Tenant Project Variable"
	TenantProjectVariableResourceName        = "tenant_project_variable"
)

type TenantProjectVariableSchema struct{}

var _ EntitySchema = TenantProjectVariableSchema{}

func (t TenantProjectVariableSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing tenants.",
		Attributes: map[string]datasourceSchema.Attribute{
			"tenant_ids":      GetQueryIDsDatasourceSchema(),
			"project_ids":     GetQueryIDsDatasourceSchema(),
			"environment_ids": GetQueryIDsDatasourceSchema(),
			"space_id":        GetSpaceIdDatasourceSchema("tenant projects", false),
			"tenant_projects": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "A list of related tenants, projects and environments that match the filter(s).",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id": GetIdDatasourceSchema(true),
						"tenant_id": datasourceSchema.StringAttribute{
							Description: "The tenant ID associated with this tenant.",
							Computed:    true,
						},
						"project_id": datasourceSchema.StringAttribute{
							Description: "The project ID associated with this tenant.",
							Computed:    true,
						},
						"environment_ids": datasourceSchema.ListAttribute{
							Description: "The environment IDs associated with this tenant.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (t TenantProjectVariableSchema) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a tenant project variable in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id": util.ResourceString().
				Computed().
				Description("The unique ID for this resource.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"space_id": util.ResourceString().
				Optional().
				Computed().
				Description("The space ID associated with this Tenant Project Variable.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"tenant_id": util.ResourceString().
				Required().
				Description("The ID of the tenant.").
				Build(),
			"project_id": util.ResourceString().
				Required().
				Description("The ID of the project.").
				Build(),
			"environment_id": util.ResourceString().
				Required().
				Description("The ID of the environment.").
				Build(),
			"template_id": util.ResourceString().
				Required().
				Description("The ID of the variable template.").
				Build(),
			"value": util.ResourceString().
				Optional().
				Sensitive().
				Description("The value of the variable.").
				Build(),
		},
	}
}
