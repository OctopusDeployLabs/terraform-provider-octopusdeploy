package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"slices"
	"strings"
)

type tagTypeResource struct {
	*Config
}

func NewTagResource() resource.Resource {
	return &tagTypeResource{}
}

func (r *tagTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TagResourceName)
}

func (r *tagTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TagSchema{}.GetResourceSchema()
}

func (r *tagTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *tagTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.ValueString() == "" || !strings.Contains(data.ID.ValueString(), "/") {
		resp.Diagnostics.AddError(`unable to import tag; ID must be "TagSets-{ID}/Tags-{ID}"`, "tag ID is null or does not match the required TagSets-{ID}/Tags-{ID} format")
		return
	}

	name := data.Name.ValueString()
	tagSetID := data.TagSetId.ValueString()
	tagSetSpaceID := data.TagSetSpaceId.ValueString()

	// if name and tag set ID are empty then an import is underway
	if name == "" && tagSetID == "" {
		tflog.Info(ctx, fmt.Sprintf("importing tag (%s)", data.ID.ValueString()))
		tagSetID = strings.Split(data.ID.ValueString(), "/")[0]
	} else {
		tflog.Info(ctx, fmt.Sprintf("reading tag (%s)", data.ID.ValueString()))
	}

	tagSet, err := tagsets.GetByID(r.Config.Client, tagSetSpaceID, tagSetID)
	if err != nil {
		processUnknownTagSetError(ctx, data, err, resp.Diagnostics)
	}

	tag := schemas.MapFromStateToTag(data)
	findByIdOrNameAndSetTag(ctx, data, tag, tagSet)

	tflog.Info(ctx, fmt.Sprintf("Tag read (%s)", tag.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tagTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagCreate(ctx, data, resp.Diagnostics, r.Client)

	tflog.Info(ctx, fmt.Sprintf("tag created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t *tagTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data, state *schemas.TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tagSetID := data.TagSetId.ValueString()
	tagSetSpaceID := data.TagSetSpaceId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("updating tag (%s)", state.ID))

	// if the tag is reassigned to another tag set
	if !data.TagSetId.Equal(state.TagSetId) {
		sourceTagSetID, destinationTagSetID := state.TagSetId.ValueString(), data.TagSetId.ValueString()
		targetSpaceId := util.Ternary(data.TagSetSpaceId.ValueString() == "", data.TagSetSpaceId, state.TagSetSpaceId)
		sourceTagSetSpaceID, destinationTagSetSpaceID := state.TagSetSpaceId.ValueString(), targetSpaceId.ValueString()

		sourceTagSet, err := tagsets.GetByID(t.Client, sourceTagSetSpaceID, sourceTagSetID)
		if err != nil {
			// if spaceID has changed, tag has been deleted, recreate required
			if !targetSpaceId.Equal(state.TagSetSpaceId) {
				tagCreate(ctx, data, resp.Diagnostics, t.Client)
				if !resp.Diagnostics.HasError() {
					resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				}
				return
			}
			util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to get source tag set", err.Error())
			return
		}

		destinationTagSet, err := tagsets.GetByID(t.Client, destinationTagSetSpaceID, destinationTagSetID)
		if err != nil {
			util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to get destination tag set", err.Error())
			return
		}

		// check to see if the name already exists in the destination tag set
		for _, tag := range destinationTagSet.Tags {
			if tag.Name == name {
				resp.Diagnostics.AddError("Tag name already exists", fmt.Sprintf("the tag name '%s' is already in use by another tag in this tag set; tag names must be unique", name))
				return
			}
		}

		tag := schemas.MapFromStateToTag(data)
		if tag.ID == "" {
			tag.ID = destinationTagSet.GetID() + "/" + strings.Split(state.ID.ValueString(), "/")[1]
		}

		// check to see that the tag is not applied to a tenant
		isUsed, err := isTagUsedByTenants(ctx, t.Client, sourceTagSetSpaceID, tag)
		if err != nil {
			data.ID = types.StringValue("")
			util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to check if tag is used by tenants", err.Error())
			return
		}

		if isUsed {
			data.ID = types.StringValue("")
			resp.Diagnostics.AddError("Tag in use", "the tag may not be transferred; it is being used by one or more tenant(s)")
			return
		}

		// all requirements are met; it is OK to transfer the tag

		// remove the tag from the source tag set and update through the API
		for i := 0; i < len(sourceTagSet.Tags); i++ {
			if sourceTagSet.Tags[i].ID == state.ID.ValueString() {
				sourceTagSet.Tags = slices.Delete(sourceTagSet.Tags, i, i+1)
				if _, err := tagsets.Update(t.Client, sourceTagSet); err != nil {
					util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to update source tag set", err.Error())
					return
				}
				break
			}
		}

		// update and add the tag to the destination tag set
		tag.ID = destinationTagSet.GetID() + "/" + strings.Split(tag.ID, "/")[1]
		destinationTagSet.Tags = append(destinationTagSet.Tags, tag)

		updatedTagSet, err := tagsets.Update(t.Client, destinationTagSet)
		if err != nil {
			util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to update destination tag set", err.Error())
			return
		}

		if err := findByIdOrNameAndSetTag(ctx, data, tag, updatedTagSet); err != nil {
			resp.Diagnostics.AddError("Failed to find updated tag", "")
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	tagSet, err := tagsets.GetByID(t.Client, tagSetSpaceID, tagSetID)
	if err != nil {
		processUnknownTagSetError(ctx, data, err, resp.Diagnostics)
		return
	}

	// find and update the tag that matches the one updated in configuration
	var updatedTag *tagsets.Tag
	for i := 0; i < len(tagSet.Tags); i++ {
		if tagSet.Tags[i].ID == state.ID.ValueString() || tagSet.Tags[i].Name == data.Name.ValueString() {
			tagSet.Tags[i] = schemas.MapFromStateToTag(data)

			updatedTagSet, err := tagsets.Update(t.Client, tagSet)
			if err != nil {
				util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "Failed to update tag set", err.Error())
				return
			}

			// Find the updated tag in the updatedTagSet
			for _, t := range updatedTagSet.Tags {
				if t.ID == state.ID.ValueString() || t.Name == data.Name.ValueString() {
					updatedTag = t
					break
				}
			}

			if updatedTag == nil {
				resp.Diagnostics.AddError("Updated tag not found", "The updated tag was not found in the response from the API")
				return
			}

			schemas.MapFromTagToState(data, updatedTag, updatedTagSet)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
	}

	resp.Diagnostics.AddError("Unable to update tag", "Tag not found in tag set")
}

func (r *tagTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data *schemas.TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagSetID := data.TagSetId.ValueString()
	tagSetSpaceID := data.TagSetSpaceId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("deleting tag (%s)", data.ID))

	tagSet, err := tagsets.GetByID(r.Config.Client, tagSetSpaceID, tagSetID)
	if err != nil {
		processUnknownTagSetError(ctx, data, err, resp.Diagnostics)
		return
	}

	tag := schemas.MapFromStateToTag(data)

	// verify tag is not associated with a tenant
	isUsed, err := isTagUsedByTenants(ctx, r.Config.Client, tagSetSpaceID, tag)
	if err != nil {
		data.ID = types.StringValue("")
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to check if tag is used by tenants", err.Error())
	}

	if isUsed {
		data.ID = types.StringValue("")
		resp.Diagnostics.AddError("the tag may not be deleted; it is being used by one or more tenant(s)", "")
	}

	// tag is known and not associated with a tenant, therefore it may be deleted

	for i := 0; i < len(tagSet.Tags); i++ {
		if tagSet.Tags[i].ID == data.ID.ValueString() {
			tagSet.Tags = slices.Delete(tagSet.Tags, i, i+1)

			if _, err := tagsets.Update(r.Config.Client, tagSet); err != nil {
				util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to update tag", err.Error())
				return
			}

			log.Printf("[INFO] tag deleted (%s)", data.ID.ValueString())
			data.ID = types.StringValue("")
			return
		}
	}
}

func tagCreate(ctx context.Context, data *schemas.TagResourceModel, diag diag.Diagnostics, client *client.Client) diag.Diagnostics {
	tflog.Info(ctx, "creating tag")

	tagSetID := data.TagSetId.ValueString()
	tagSetSpaceID := data.TagSetSpaceId.ValueString()

	tagSet, err := tagsets.GetByID(client, tagSetSpaceID, tagSetID)

	if err != nil {
		processUnknownTagSetError(ctx, data, err, diag)
		return diag
	}

	name := data.Name.ValueString()

	for _, tag := range tagSet.Tags {
		if tag.Name == name {
			diag.AddError(`the tag name '%s' is already in use by another tag in this tag set; tag names must be unique`, name)
		}
	}

	tag := schemas.MapFromStateToTag(data)
	if tag.ID != "" {
		tag.ID = tagSet.GetID() + "/" + strings.Split(tag.ID, "/")[1]
	}
	tagSet.Tags = append(tagSet.Tags, tag)

	updatedTagSet, err := tagsets.Update(client, tagSet)
	if err != nil {
		diag.AddError(`unable to update tag set`, err.Error())
	}

	return findByIdOrNameAndSetTag(ctx, data, tag, updatedTagSet)
}

func isTagUsedByTenants(ctx context.Context, octopus *client.Client, spaceID string, tag *tagsets.Tag) (bool, error) {
	tenants, err := tenants.Get(octopus, spaceID, tenants.TenantsQuery{
		Tags: []string{tag.ID},
	})
	if err != nil {
		return false, err
	}

	return len(tenants.Items) > 0, nil
}

func findByIdOrNameAndSetTag(ctx context.Context, data *schemas.TagResourceModel, tag *tagsets.Tag, tagSet *tagsets.TagSet) diag.Diagnostics {
	for _, t := range tagSet.Tags {
		if t.Name == tag.Name {
			schemas.MapFromTagToState(data, t, tagSet)

			tflog.Info(ctx, fmt.Sprintf("tag (%s)", tag.ID))
			return nil
		}
	}

	for _, t := range tagSet.Tags {
		if t.ID == tag.ID {
			schemas.MapFromTagToState(data, t, tagSet)
			tflog.Info(ctx, fmt.Sprintf("tag (%s)", tag.ID))
			return nil
		}
	}

	tflog.Info(ctx, fmt.Sprintf("%s (%s) not found; deleting from state", tag.ID, tag.ID))
	data.ID = types.StringValue("")
	return nil
}

func processUnknownTagSetError(ctx context.Context, data *schemas.TagResourceModel, err error, diag diag.Diagnostics) {
	if err == nil {
		return
	}

	if apiError, ok := err.(*core.APIError); ok {
		if apiError.StatusCode == 404 {
			tflog.Info(ctx, fmt.Sprintf("tag set (%s) not found; deleting tag from state", data.ID))
			data.ID = types.StringValue("")
			return
		}
	}

	diag.AddError("Processing unknown tag set failed", err.Error())
}

func (t *tagTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
