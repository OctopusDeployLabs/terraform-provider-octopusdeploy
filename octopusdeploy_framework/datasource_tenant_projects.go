package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

type TenantProjectsDataModel struct {
	SpaceID        types.String `tfsdk:"space_id"`
	TenantIDs      types.List   `tfsdk:"tenant_ids"`
	ProjectIDs     types.List   `tfsdk:"project_ids"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
	TenantProjects types.List   `tfsdk:"tenant_projects"`
}

type tenantProjectsDataSource struct {
	*Config
}

func NewTenantProjectDataSource() datasource.DataSource {
	return &tenantProjectsDataSource{}
}

func (*tenantProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("tenant_projects")
}

func (t *tenantProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	t.Config = DataSourceConfiguration(req, resp)
}

func (*tenantProjectsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Provides information about existing tenants.",
		Attributes: map[string]datasourceSchema.Attribute{
			"tenant_ids":      schemas.GetQueryIDsDatasourceSchema(),
			"project_ids":     schemas.GetQueryIDsDatasourceSchema(),
			"environment_ids": schemas.GetQueryIDsDatasourceSchema(),
			"space_id":        schemas.GetSpaceIdDatasourceSchema("tenant projects", false),
		},
		Blocks: map[string]datasourceSchema.Block{
			"tenant_projects": datasourceSchema.ListNestedBlock{
				Description: "A list of related tenants, projects and environments that match the filter(s).",
				NestedObject: datasourceSchema.NestedBlockObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id": schemas.GetIdDatasourceSchema(true),
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

func (t *tenantProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TenantProjectsDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantIDs := util.ExpandStringList(data.TenantIDs)
	projectIDs := util.ExpandStringList(data.ProjectIDs)
	environmentIDs := util.ExpandStringList(data.EnvironmentIDs)
	spaceID := data.SpaceID.ValueString()

	if (tenantIDs == nil || len(tenantIDs) == 0) && (projectIDs == nil || len(projectIDs) == 0) {
		resp.Diagnostics.AddError("must provide at least one tenant or one project", "tenant IDs and project IDs are nil")
		return
	}

	var tenantData []*tenants.Tenant
	if tenantIDs == nil {
		tenantData = make([]*tenants.Tenant, 0)
		for _, projectID := range projectIDs {
			// todo: should consider using concurrency
			tenantsByProjectID, err := getTenantByProjectID(t.Client, projectID, spaceID)
			if err != nil {
				resp.Diagnostics.AddError("unable to load tenant data for project", err.Error())
				return
			}
			tenantData = append(tenantData, tenantsByProjectID...)
		}
	} else {
		tenantQuery := tenants.TenantsQuery{
			IDs: tenantIDs,
		}
		tenantResource, err := tenants.Get(t.Client, spaceID, tenantQuery)
		if err != nil {
			resp.Diagnostics.AddError("unable to load tenant data", err.Error())
			return
		}
		tenantData, err = tenantResource.GetAllPages(t.Client.Sling())
		if err != nil {
			resp.Diagnostics.AddError("unable to load tenant data", err.Error())
			return
		}
	}

	// todo: refactor
	flattenedTenantProjects := make([]any, 0)
	for _, tenant := range tenantData {
		if projectIDs == nil {
			for projectID, envIDs := range tenant.ProjectEnvironments {
				if environmentIDs != nil {
					for _, envID := range envIDs {
						if slices.Contains(environmentIDs, envID) {
							flattenedTenantProjects = append(flattenedTenantProjects, flattenTenantProject(tenant, projectID))
							break
						}
					}
				} else {
					flattenedTenantProjects = append(flattenedTenantProjects, flattenTenantProject(tenant, projectID))
				}
			}
			continue
		}
		for _, projectID := range projectIDs {
			if envIDs, ok := tenant.ProjectEnvironments[projectID]; ok {
				if environmentIDs != nil {
					for _, envID := range envIDs {
						if slices.Contains(environmentIDs, envID) {
							flattenedTenantProjects = append(flattenedTenantProjects, flattenTenantProject(tenant, projectID))
							break
						}
					}
				} else {
					flattenedTenantProjects = append(flattenedTenantProjects, flattenTenantProject(tenant, projectID))
				}
			}

		}
	}

	data.TenantProjects, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: tenantProjectType()}, flattenedTenantProjects)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getTenantByProjectID(c *client.Client, projectID string, spaceID string) ([]*tenants.Tenant, error) {
	tenantQuery := tenants.TenantsQuery{
		ProjectID: projectID,
	}
	tenantResource, err := tenants.Get(c, spaceID, tenantQuery)
	if err != nil {
		return nil, err
	}
	tenantData, err := tenantResource.GetAllPages(c.Sling())
	if err != nil {
		return nil, err
	}
	return tenantData, nil
}

func buildTenantProjectID(spaceID string, tenantID string, projectID string) string {
	return fmt.Sprintf("%s:%s:%s", spaceID, tenantID, projectID)
}

func tenantProjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"tenant_id":       types.StringType,
		"project_id":      types.StringType,
		"environment_ids": types.ListType{ElemType: types.StringType},
	}
}

func flattenTenantProject(tenant *tenants.Tenant, projectID string) attr.Value {
	environmentIDs := make([]attr.Value, len(tenant.ProjectEnvironments[projectID]))
	for i, envID := range tenant.ProjectEnvironments[projectID] {
		environmentIDs[i] = types.StringValue(envID)
	}

	environmentIdList, _ := types.ListValue(types.StringType, environmentIDs)

	return types.ObjectValueMust(tenantProjectType(), map[string]attr.Value{
		"id":              types.StringValue(buildTenantProjectID(tenant.SpaceID, tenant.ID, projectID)),
		"tenant_id":       types.StringValue(tenant.ID),
		"project_id":      types.StringValue(projectID),
		"environment_ids": environmentIdList,
	})
}
