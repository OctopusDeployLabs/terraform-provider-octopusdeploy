package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ociRegistryFeedTypeResource struct {
	*Config
}

func NewOCIRegistryFeedResource() resource.Resource {
	return &ociRegistryFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &ociRegistryFeedTypeResource{}

func (r *ociRegistryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("oci_registry_feed")
}

func (r *ociRegistryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.OCIRegistryFeedSchema{}.GetResourceSchema()
}

func (r *ociRegistryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *ociRegistryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.OCIRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	feed, err := createOCIRegistryResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating OCI Registry feed: %s", feed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, feed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create OCI Registry feed", err.Error())
		return
	}

	updateDataFromOCIRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.OCIRegistryFeed))

	tflog.Info(ctx, fmt.Sprintf("OCI Registry feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ociRegistryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.OCIRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading OCI Registry feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "OCI Registry feed"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load OCI Registry feed", err.Error())
		}
		return
	}

	loadedFeed := feed.(*feeds.OCIRegistryFeed)
	updateDataFromOCIRegistryFeed(data, data.SpaceID.ValueString(), loadedFeed)

	tflog.Info(ctx, fmt.Sprintf("OCI Registry feed read (%s)", loadedFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ociRegistryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.OCIRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating OCI Registry feed '%s'", data.ID.ValueString()))

	feed, err := createOCIRegistryResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load OCI Registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating OCI Registry feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update OCI Registry feed", err.Error())
		return
	}

	updateDataFromOCIRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.OCIRegistryFeed))

	tflog.Info(ctx, fmt.Sprintf("OCI Registry feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ociRegistryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.OCIRegistryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete OCI Registry feed", err.Error())
		return
	}
}

func createOCIRegistryResourceFromData(data *schemas.OCIRegistryFeedTypeResourceModel) (*feeds.OCIRegistryFeed, error) {
	feed, err := feeds.NewOCIRegistryFeed(data.Name.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.FeedURI = data.FeedUri.ValueString()

	feed.Username = data.Username.ValueString()
	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()

	return feed, nil
}

func updateDataFromOCIRegistryFeed(data *schemas.OCIRegistryFeedTypeResourceModel, spaceId string, feed *feeds.OCIRegistryFeed) {
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}

	data.ID = types.StringValue(feed.ID)
}

func (*ociRegistryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
