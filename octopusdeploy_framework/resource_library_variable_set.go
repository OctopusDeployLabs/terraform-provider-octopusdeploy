package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
)

type libraryVariableSetFeedTypeResource struct {
	*Config
}

func (r *libraryVariableSetFeedTypeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *schemas.LibraryVariableSetResourceModel

	if req.Plan.Raw.IsNull() {
		return
	}

	//if req.State.Raw.IsNull() {
	//	isCreation = true
	//} else {
	//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	//}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//templates := plan.Template
	//expandedActionTemplates := schemas.ExpandActionTemplateParameters(templates)
	//templateIdsValues := schemas.FlattenTemplateIds(expandedActionTemplates)
	//resp.Plan.SetAttribute(ctx, path.Root("template_ids"), templateIdsValues)
}

func NewLibraryVariableSetFeedResource() resource.Resource {
	return &libraryVariableSetFeedTypeResource{}
}

func (r *libraryVariableSetFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("library_variable_set")
}

func (r *libraryVariableSetFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetLibraryVariableSetResourceSchema()
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
		resp.Diagnostics.AddError("unable to create library variable set", err.Error())
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
		resp.Diagnostics.AddError("unable to load library variable set", err.Error())
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
		resp.Diagnostics.AddError("unable to update library variable set", err.Error())
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
		resp.Diagnostics.AddError("unable to delete library variable set", err.Error())
		return
	}
}
