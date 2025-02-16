package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type libraryVariableSetFeedTypeResource struct {
	*Config
}

func NewLibraryVariableSetFeedResource() resource.Resource {
	return &libraryVariableSetFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &libraryVariableSetFeedTypeResource{}

func (r *libraryVariableSetFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("library_variable_set")
}

func (r *libraryVariableSetFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.LibraryVariableSetSchema{}.GetResourceSchema()
}

func (r *libraryVariableSetFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *libraryVariableSetFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.LibraryVariableSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newLibraryVariableSet := schemas.MapToLibraryVariableSet(data)
	libraryVariableSet, err := libraryvariablesets.Add(r.Config.Client, newLibraryVariableSet)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "unable to create library variable set", err.Error())
		return
	}

	schemas.MapFromLibraryVariableSet(data, libraryVariableSet.SpaceID, libraryVariableSet)
	tflog.Info(ctx, fmt.Sprintf("Library Variable Set created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *libraryVariableSetFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.LibraryVariableSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[INFO] reading library variable set (%s)", data.ID.ValueString())

	libraryVariableSet, err := libraryvariablesets.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "library variable set"); err != nil {
			util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "unable to load library variable set", err.Error())
		}
		return
	}

	schemas.MapFromLibraryVariableSet(data, data.SpaceID.ValueString(), libraryVariableSet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *libraryVariableSetFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.LibraryVariableSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating library variable set '%s'", data.ID.ValueString()))

	libraryVariableSet := schemas.MapToLibraryVariableSet(data)
	libraryVariableSet.ID = state.ID.ValueString()

	updatedLibraryVariableSet, err := libraryvariablesets.Update(r.Config.Client, libraryVariableSet)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "unable to update library variable set", err.Error())
		return
	}
	schemas.MapFromLibraryVariableSet(data, state.SpaceID.ValueString(), updatedLibraryVariableSet)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *libraryVariableSetFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.LibraryVariableSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := libraryvariablesets.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "unable to delete library variable set", err.Error())
		return
	}
}

func (*libraryVariableSetFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
