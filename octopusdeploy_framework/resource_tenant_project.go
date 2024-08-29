package octopusdeploy_framework

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tenantProjectResource struct {
	*Config
}

var _ resource.Resource = &tenantProjectResource{}
var _ resource.ResourceWithImportState = &tenantProjectResource{}
var _ resource.ResourceWithConfigure = &tenantProjectResource{}

func NewTenantProjectResource() resource.Resource {
	return &tenantProjectResource{}
}

func (t *tenantProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("tenant_project")
}

func (t *tenantProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schemas.GetIdResourceSchema(),
			"tenant_id": schema.StringAttribute{
				Description: "The tenant ID associated with this tenant.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID associated with this tenant.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_ids": schema.ListAttribute{
				Description: "The environment IDs associated with this tenant.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"space_id": schemas.GetSpaceIdResourceSchema("project tenant"),
		}}
	resp.Schema = schemas.GetTenantProjectsResourceSchema()
}

func (t *tenantProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tenantProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan schemas.TenantProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := t.getSpaceId(plan)

	tflog.Info(ctx, fmt.Sprintf("connecting tenant (%s) to project (%s)", plan.TenantID, plan.ProjectID))

	tenant, err := tenants.GetByID(t.Client, spaceId, plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot load tenant", err.Error())
		return
	}

	tenant.ProjectEnvironments[plan.ProjectID.ValueString()] = util.ExpandStringList(plan.EnvironmentIDs)

	_, err = tenants.Update(t.Client, tenant)
	if err != nil {
		resp.Diagnostics.AddError("cannot update tenant environment", err.Error())
	}

	plan.ID = types.StringValue(schemas.BuildTenantProjectID(spaceId, plan.TenantID.ValueString(), plan.ProjectID.ValueString()))
	plan.SpaceID = types.StringValue(spaceId)
	plan.EnvironmentIDs = util.FlattenStringList(tenant.ProjectEnvironments[plan.ProjectID.ValueString()])

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, fmt.Sprintf("tenant (%s) connected to project (%#v)", plan.TenantID.ValueString(), plan.ProjectID.ValueString()))
}

func (t *tenantProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.TenantProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bits := strings.Split(data.ID.ValueString(), ":")
	spaceID := bits[0]
	tenantID := bits[1]
	projectID := bits[2]

	tenant, err := tenants.GetByID(t.Client, spaceID, tenantID)
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError("unable to load tenant", err.Error())
			return
		}
	}

	data.EnvironmentIDs = util.FlattenStringList(tenant.ProjectEnvironments[projectID])
	data.SpaceID = types.StringValue(spaceID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t *tenantProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	// read plan and state
	var plan, state schemas.TenantProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := t.getSpaceId(plan)

	tenant, err := tenants.GetByID(t.Client, spaceId, plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot load tenant", err.Error())
		return
	}

	tenant.ProjectEnvironments[plan.ProjectID.ValueString()] = util.ExpandStringList(plan.EnvironmentIDs)

	_, err = tenants.Update(t.Client, tenant)
	if err != nil {
		resp.Diagnostics.AddError("cannot update tenant environment", err.Error())
	}

	plan.ID = types.StringValue(schemas.BuildTenantProjectID(spaceId, plan.TenantID.ValueString(), plan.ProjectID.ValueString()))
	plan.SpaceID = types.StringValue(spaceId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, fmt.Sprintf("updated tenant (%s) connection to project (%#v)", plan.TenantID.ValueString(), plan.ProjectID.ValueString()))
}

func (t *tenantProjectResource) getSpaceId(plan schemas.TenantProjectResourceModel) string {
	spaceId := plan.SpaceID.ValueString()
	if spaceId == "" {
		spaceId = t.Client.GetSpaceID()
	}
	return spaceId
}

func (t *tenantProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()
	var data schemas.TenantProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("removing tenant (%s) from project (%s)", data.TenantID.ValueString(), data.ProjectID.ValueString()))

	spaceId := t.getSpaceId(data)

	tenant, err := tenants.GetByID(t.Client, spaceId, data.TenantID.ValueString())
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode == http.StatusNotFound {
			tflog.Info(ctx, fmt.Sprintf("tenant (%s) no longer exists", data.TenantID.ValueString()))
			return
		} else {
			resp.Diagnostics.AddError("cannot load tenant", err.Error())
			return
		}
	}

	delete(tenant.ProjectEnvironments, data.ProjectID.ValueString())
	_, err = tenants.Update(t.Client, tenant)
	if err != nil {
		resp.Diagnostics.AddError("cannot remove tenant environment", err.Error())
	}

	tflog.Info(ctx, fmt.Sprintf("tenant (%s) disconnected from project (%s)", data.TenantID.ValueString(), data.ProjectID.ValueString()))
}

func (t *tenantProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
