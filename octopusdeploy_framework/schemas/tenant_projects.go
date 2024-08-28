package schemas

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantProjectsDataModel struct {
	SpaceID        types.String `tfsdk:"space_id"`
	TenantIDs      types.List   `tfsdk:"tenant_ids"`
	ProjectIDs     types.List   `tfsdk:"project_ids"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
	TenantProjects types.List   `tfsdk:"tenant_projects"`
}

type TenantProjectResourceModel struct {
	SpaceID        types.String `tfsdk:"space_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`

	ResourceModel
}

func GetTenantProjectsDataSourceSchema() datasourceSchema.Schema {
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

func GetTenantProjectsResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id": util.GetIdResourceSchema(),
			"tenant_id": resourceSchema.StringAttribute{
				Description: "The tenant ID associated with this tenant.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": resourceSchema.StringAttribute{
				Description: "The project ID associated with this tenant.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_ids": resourceSchema.ListAttribute{
				Description: "The environment IDs associated with this tenant.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"space_id": GetSpaceIdResourceSchema("project tenant"),
		}}
}

func BuildTenantProjectID(spaceID string, tenantID string, projectID string) string {
	return fmt.Sprintf("%s:%s:%s", spaceID, tenantID, projectID)
}

func TenantProjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"tenant_id":       types.StringType,
		"project_id":      types.StringType,
		"environment_ids": types.ListType{ElemType: types.StringType},
	}
}

func MapTenantToTenantProject(tenant *tenants.Tenant, projectID string) attr.Value {
	environmentIDs := make([]attr.Value, len(tenant.ProjectEnvironments[projectID]))
	for i, envID := range tenant.ProjectEnvironments[projectID] {
		environmentIDs[i] = types.StringValue(envID)
	}

	environmentIdList, _ := types.ListValue(types.StringType, environmentIDs)

	return types.ObjectValueMust(TenantProjectType(), map[string]attr.Value{
		"id":              types.StringValue(BuildTenantProjectID(tenant.SpaceID, tenant.ID, projectID)),
		"tenant_id":       types.StringValue(tenant.ID),
		"project_id":      types.StringValue(projectID),
		"environment_ids": environmentIdList,
	})
}
