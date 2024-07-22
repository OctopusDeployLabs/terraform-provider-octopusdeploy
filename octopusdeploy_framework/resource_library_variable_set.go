package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

//TODO: Plan Modifiers for
//		CustomizeDiff: fixTemplateIds,

type libraryVariableSetFeedTypeResource struct {
	*Config
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

// fixTemplateIds uses the suggestion from https://github.com/hashicorp/terraform/issues/18863
// to ensure that the template_ids field has keys to match the list of template names.
func fixTemplateIds(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	templates := d.Get("template")
	templateIds := map[string]string{}
	if templates != nil {
		for _, t := range templates.([]interface{}) {
			template := t.(map[string]interface{})
			templateIds[template["name"].(string)] = template["id"].(string)
		}
	}
	if err := d.SetNew("template_ids", templateIds); err != nil {
		return err
	}

	return nil
}

func (r *libraryVariableSetFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.LibraryVariableSetResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newLibraryVariableSet := schemas.CreateLibraryVariableSet(data)
	libraryVariableSet, err := libraryvariablesets.Add(r.Config.Client, newLibraryVariableSet)
	if err != nil {
		resp.Diagnostics.AddError("unable to create library variable set", err.Error())
		return
	}

	schemas.UpdateDataFromLibraryVariableSet(data, data.SpaceID.ValueString(), libraryVariableSet)
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
	schemas.UpdateDataFromLibraryVariableSet(data, data.SpaceID.ValueString(), updatedLibraryVariableSet)
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
