package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"reflect"
)

//TODO: Plan Modifiers for
//		CustomizeDiff: fixTemplateIds,

type libraryVariableSetFeedTypeResource struct {
	*Config
}

type libraryVariableSetFeedModifier struct{}

func (m libraryVariableSetFeedModifier) Description(_ context.Context) string {
	return "Template_Ids will be populated once created"
}

func (m libraryVariableSetFeedModifier) MarkdownDescription(_ context.Context) string {
	return "Template_Ids will be populated once created"
}

func (m libraryVariableSetFeedModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}

	oldValues := req.StateValue.Elements()
	for k, v := range req.PlanValue.Elements() {
		o, ok := oldValues[k]
		log.Println(o)
		if !ok {
			// something new, go ahead with the plan
			return
		}

		s, ok := v.(types.String)
		log.Println(s)
		if !ok {
			// this is very bad and shouldn't haven't gotten past validation
			resp.Diagnostics.AddError(fmt.Sprintf("Invalid value type for json string in plan for %s", req.Path), "invalid value")
			return
		}

		var newValue map[string]interface{}

		s, ok = o.(types.String)

		var oldValue map[string]interface{}
		//err = json.Unmarshal([]byte(s.ValueString()), &oldValue)
		//if err != nil {
		//	resp.Diagnostics.AddError(fmt.Sprintf("Invalid json in plan for %s", req.Path), err.Error())
		//	return
		//}

		if !reflect.DeepEqual(oldValue, newValue) {
			return
		}

		delete(oldValues, k)
	}

	if len(oldValues) > 0 {
		// something was removed in the plan
		return
	}

	resp.PlanValue = req.StateValue
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

	newLibraryVariableSet := schemas.CreateLibraryVariableSet(data)
	libraryVariableSet, err := libraryvariablesets.Add(r.Config.Client, newLibraryVariableSet)
	if err != nil {
		resp.Diagnostics.AddError("unable to create library variable set", err.Error())
		return
	}

	schemas.UpdateDataFromLibraryVariableSet(data, libraryVariableSet.SpaceID, libraryVariableSet)
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

	schemas.UpdateDataFromLibraryVariableSet(data, data.SpaceID.ValueString(), libraryVariableSet)

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

	libraryVariableSet := schemas.CreateLibraryVariableSet(data)
	libraryVariableSet.ID = state.ID.ValueString()

	updatedLibraryVariableSet, err := libraryvariablesets.Update(r.Config.Client, libraryVariableSet)
	if err != nil {
		resp.Diagnostics.AddError("unable to update library variable set", err.Error())
		return
	}
	schemas.UpdateDataFromLibraryVariableSet(data, state.SpaceID.ValueString(), updatedLibraryVariableSet)
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
