package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type artifactoryGenericFeedTypeResource struct {
	*Config
}

func NewArtifactoryGenericFeedResource() resource.Resource {
	return &artifactoryGenericFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &artifactoryGenericFeedTypeResource{}

func (r *artifactoryGenericFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("artifactory_generic_feed")
}

func (r *artifactoryGenericFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ArtifactoryGenericFeedSchema{}.GetResourceSchema()
}

func (r *artifactoryGenericFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *artifactoryGenericFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ArtifactoryGenericFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	artifactoryGenericFeed, err := createArtifactoryGenericResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating ArtifactoryGeneric feed: %s", artifactoryGenericFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, artifactoryGenericFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create artifactoryGeneric feed", err.Error())
		return
	}

	updateDataFromArtifactoryGenericFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.ArtifactoryGenericFeed))

	tflog.Info(ctx, fmt.Sprintf("ArtifactoryGeneric feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *artifactoryGenericFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ArtifactoryGenericFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading ArtifactoryGeneric feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "artifactory generic feed"); err != nil {
			resp.Diagnostics.AddError("unable to load artifactoryGeneric feed", err.Error())
		}
		return
	}

	artifactoryGenericFeed := feed.(*feeds.ArtifactoryGenericFeed)
	updateDataFromArtifactoryGenericFeed(data, data.SpaceID.ValueString(), artifactoryGenericFeed)

	tflog.Info(ctx, fmt.Sprintf("ArtifactoryGeneric feed read (%s)", artifactoryGenericFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *artifactoryGenericFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.ArtifactoryGenericFeedTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating artifactoryGeneric feed '%s'", data.ID.ValueString()))

	feed, err := createArtifactoryGenericResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load artifactoryGeneric feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating ArtifactoryGeneric feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update artifactoryGeneric feed", err.Error())
		return
	}

	updateDataFromArtifactoryGenericFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.ArtifactoryGenericFeed))

	tflog.Info(ctx, fmt.Sprintf("ArtifactoryGeneric feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *artifactoryGenericFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.ArtifactoryGenericFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete artifactoryGeneric feed", err.Error())
		return
	}
}

func createArtifactoryGenericResourceFromData(data *schemas.ArtifactoryGenericFeedTypeResourceModel) (*feeds.ArtifactoryGenericFeed, error) {
	feed, err := feeds.NewArtifactoryGenericFeed(data.Name.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.FeedURI = data.FeedUri.ValueString()

	var packageAcquisitionLocationOptions []string
	for _, element := range data.PackageAcquisitionLocationOptions.Elements() {
		packageAcquisitionLocationOptions = append(packageAcquisitionLocationOptions, element.(types.String).ValueString())
	}

	feed.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptions
	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()
	feed.Username = data.Username.ValueString()
	feed.Repository = data.Repository.ValueString()
	feed.LayoutRegex = data.LayoutRegex.ValueString()

	return feed, nil
}

func updateDataFromArtifactoryGenericFeed(data *schemas.ArtifactoryGenericFeedTypeResourceModel, spaceId string, feed *feeds.ArtifactoryGenericFeed) {
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}
	data.Repository = types.StringValue(feed.Repository)
	if feed.LayoutRegex != "" {
		data.LayoutRegex = types.StringValue(feed.LayoutRegex)
	}

	packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
	for i, option := range feed.PackageAcquisitionLocationOptions {
		packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
	}

	var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
	data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	data.ID = types.StringValue(feed.GetID())
}

func (*artifactoryGenericFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
