package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"slices"
	"sync"
)

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
	resp.Schema = schemas.TenantProjectVariableSchema{}.GetDatasourceSchema()
}

func (t *tenantProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.TenantProjectsDataModel
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

	tenantData, err := getTenantData(t.Client, tenantIDs, projectIDs, spaceID)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, t.Config.SystemInfo, "unable to load tenant data", err.Error())
		return
	}

	mappedTenantProjects := make([]any, 0)

	filterEnvAndMap := func(tenant *tenants.Tenant, projectID string, envIDs []string) {
		if environmentIDs == nil || len(environmentIDs) == 0 {
			mappedTenantProjects = append(mappedTenantProjects, schemas.MapTenantToTenantProject(tenant, projectID))
			return
		}
		for _, envID := range envIDs {
			if slices.Contains(environmentIDs, envID) {
				mappedTenantProjects = append(mappedTenantProjects, schemas.MapTenantToTenantProject(tenant, projectID))
				return
			}
		}
	}

	if projectIDs == nil || len(projectIDs) == 0 {
		for _, tenant := range tenantData {
			for projectID, envIDs := range tenant.ProjectEnvironments {
				filterEnvAndMap(tenant, projectID, envIDs)
			}
		}
	} else {
		for _, tenant := range tenantData {
			for _, projectID := range projectIDs {
				if envIDs, ok := tenant.ProjectEnvironments[projectID]; ok {
					filterEnvAndMap(tenant, projectID, envIDs)
				}
			}
		}
	}

	data.TenantProjects, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.TenantProjectType()}, mappedTenantProjects)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getTenantData(c *client.Client, tenantIDs, ProjectIDs []string, spaceID string) ([]*tenants.Tenant, error) {
	if tenantIDs == nil {
		return getTenantByProjectIDs(c, ProjectIDs, spaceID)
	}
	return getTenantsByTenantIDs(c, tenantIDs, spaceID)
}

const maxConcurrentTenantByProjectIDRequests = 5

func getTenantByProjectIDs(c *client.Client, projectIDs []string, spaceID string) ([]*tenants.Tenant, error) {
	tenantData := make([]*tenants.Tenant, 0)
	tenantOutCh := make(chan []*tenants.Tenant, len(projectIDs))
	errorCh := make(chan error, 1)
	var wg sync.WaitGroup
	guardCh := make(chan struct{}, maxConcurrentTenantByProjectIDRequests)

	for _, projectID := range projectIDs {
		wg.Add(1)
		guardCh <- struct{}{}
		go func(projectID string) {
			defer wg.Done()
			defer func() { <-guardCh }()
			tenantsByProjectID, err := getTenantByProjectID(c, projectID, spaceID)
			if err != nil {
				select {
				case errorCh <- fmt.Errorf("unable to load tenant data for project %s: %w", projectID, err):
				default:
					// Avoid blocking if errorCh already has value
				}
				return
			}
			tenantOutCh <- tenantsByProjectID
		}(projectID)
	}

	go func() {
		wg.Wait()
		close(tenantOutCh)
		close(errorCh)
	}()

	if err := <-errorCh; err != nil {
		return nil, err
	}

	seenTenantIDs := make(map[string]struct{})
	for tenantsByProjectID := range tenantOutCh {
		for _, tenantByProjectID := range tenantsByProjectID {
			if _, ok := seenTenantIDs[tenantByProjectID.ID]; ok {
				continue
			}
			seenTenantIDs[tenantByProjectID.ID] = struct{}{}
			tenantData = append(tenantData, tenantByProjectID)
		}
	}

	return tenantData, nil
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

func getTenantsByTenantIDs(c *client.Client, tenantIDs []string, spaceID string) ([]*tenants.Tenant, error) {
	tenantQuery := tenants.TenantsQuery{
		IDs: tenantIDs,
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
