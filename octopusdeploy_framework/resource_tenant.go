package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tenantTypeResource struct {
	*Config
}

func NewTenantResource() resource.Resource {
	return &tenantTypeResource{}
}

var _ resource.ResourceWithImportState = &tenantTypeResource{}

func (r *tenantTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("tenant")
}

func (r *tenantTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TenantSchema{}.GetResourceSchema()
}

func (r *tenantTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *tenantTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.TenantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, err := mapStateToTenant(ctx, data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Tenant: %s", tenant.Name))

	createdTenant, err := tenants.Add(r.Config.Client, tenant)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create tenant", err.Error())
		return
	}

	mapTenantToState(ctx, data, createdTenant)

	tflog.Info(ctx, fmt.Sprintf("Tenant created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tenantTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.TenantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Tenant (%s)", data.ID))

	client := r.Config.Client
	tenant, err := tenants.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "tenant"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load tenant", err.Error())
		}
		return
	}

	mapTenantToState(ctx, data, tenant)

	tflog.Info(ctx, fmt.Sprintf("Tenant read (%s)", tenant.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tenantTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data, state *schemas.TenantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating tenant '%s'", data.ID.ValueString()))

	tenantFromApi, err := tenants.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())

	tenant, err := mapStateToTenant(ctx, data)
	tenant.ID = state.ID.ValueString()
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to map to tenant", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Tenant (%s)", data.ID))

	tenant.ProjectEnvironments = tenantFromApi.ProjectEnvironments
	updatedTenant, err := tenants.Update(r.Config.Client, tenant)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update tenant", err.Error())
		return
	}

	mapTenantToState(ctx, data, updatedTenant)

	tflog.Info(ctx, fmt.Sprintf("Tenant updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tenantTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data schemas.TenantModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := tenants.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete tenant", err.Error())
		return
	}
}

func mapStateToTenant(ctx context.Context, data *schemas.TenantModel) (*tenants.Tenant, error) {
	tenant := tenants.NewTenant(data.Name.ValueString())
	tenant.ID = data.ID.ValueString()
	tenant.ClonedFromTenantID = data.ClonedFromTenantId.ValueString()
	tenant.Description = data.Description.ValueString()
	tenant.IsDisabled = data.IsDisabled.ValueBool()
	tenant.SpaceID = data.SpaceID.ValueString()

	convertedTenantTags, diags := util.SetToStringArray(ctx, data.TenantTags)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error converting tenant tags: %v\n", diags))
	}

	tenant.TenantTags = convertedTenantTags

	return tenant, nil
}

func mapTenantToState(ctx context.Context, data *schemas.TenantModel, tenant *tenants.Tenant) {
	data.ID = types.StringValue(tenant.ID)
	data.ClonedFromTenantId = types.StringValue(tenant.ClonedFromTenantID)
	data.Description = types.StringValue(tenant.Description)
	data.IsDisabled = types.BoolValue(tenant.IsDisabled)
	data.SpaceID = types.StringValue(tenant.SpaceID)
	data.Name = types.StringValue(tenant.Name)

	convertedTenantTags, diags := util.SetToStringArray(ctx, data.TenantTags)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error converting tenant tags: %v\n", diags))
	}

	data.TenantTags = basetypes.SetValue(util.FlattenStringList(convertedTenantTags))
}

func (*tenantTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
