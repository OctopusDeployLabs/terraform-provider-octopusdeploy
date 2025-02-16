package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/scriptmodules"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type scriptModuleTypeResource struct {
	*Config
}

func NewScriptModuleResource() resource.Resource {
	return &scriptModuleTypeResource{}
}

var _ resource.ResourceWithImportState = &scriptModuleTypeResource{}

func (r *scriptModuleTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("script_module")
}

func (r *scriptModuleTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ScriptModuleSchema{}.GetResourceSchema()
}

func (r *scriptModuleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *scriptModuleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.ScriptModuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptModule := schemas.MapFromScriptModuleToState(data)

	tflog.Info(ctx, fmt.Sprintf("creating Script Module: %s", scriptModule.Name))

	createdScriptModule, err := scriptmodules.Add(r.Config.Client, scriptModule)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create script module", err.Error())
		return
	}

	schemas.MapToScriptModuleFromState(data, createdScriptModule)

	tflog.Info(ctx, fmt.Sprintf("Script Module created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scriptModuleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ScriptModuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Script Module (%s)", data.ID))

	client := r.Config.Client
	scriptModule, err := scriptmodules.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "Script Module"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load script module", err.Error())
		}
		return
	}

	schemas.MapToScriptModuleFromState(data, scriptModule)

	tflog.Info(ctx, fmt.Sprintf("Script Module read (%s)", scriptModule.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scriptModuleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.ScriptModuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating script module '%s'", data.ID.ValueString()))

	scriptModule := schemas.MapFromScriptModuleToState(data)
	scriptModule.ID = state.ID.ValueString()

	updatedScriptModule, err := scriptmodules.Update(r.Config.Client, scriptModule)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update script module", err.Error())
		return
	}

	schemas.MapToScriptModuleFromState(data, updatedScriptModule)

	tflog.Info(ctx, fmt.Sprintf("Script Module updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scriptModuleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data schemas.ScriptModuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := scriptmodules.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete script module", err.Error())
		return
	}
}

func (*scriptModuleTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
