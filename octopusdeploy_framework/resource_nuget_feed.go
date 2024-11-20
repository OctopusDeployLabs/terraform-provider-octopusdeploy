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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type nugetFeedTypeResource struct {
	*Config
}

func NewNugetFeedResource() resource.Resource {
	return &nugetFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &nugetFeedTypeResource{}

func (r *nugetFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("nuget_feed")
}

func (r *nugetFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.NugetFeedSchema{}.GetResourceSchema()
}

func (r *nugetFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *nugetFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.NugetFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nugetFeed, err := createNugetResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Nuget feed: %s", nugetFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, nugetFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create nuget feed", err.Error())
		return
	}

	updateDataFromNugetFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.NuGetFeed))

	data.ID = types.StringValue(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("Nuget feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *nugetFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.NugetFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Nuget feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, &resp.State, data, err, "nuget feed"); err != nil {
			resp.Diagnostics.AddError("unable to load nuget feed", err.Error())
		}
		return
	}

	nugetFeed := feed.(*feeds.NuGetFeed)
	updateDataFromNugetFeed(data, data.SpaceID.ValueString(), nugetFeed)

	tflog.Info(ctx, fmt.Sprintf("Nuget feed read (%s)", nugetFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *nugetFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.NugetFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating nuget feed '%s'", data.ID.ValueString()))

	feed, err := createNugetResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load nuget feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Nuget feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update nuget feed", err.Error())
		return
	}

	updateDataFromNugetFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.NuGetFeed))

	tflog.Info(ctx, fmt.Sprintf("Nuget feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *nugetFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.NugetFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete nuget feed", err.Error())
		return
	}
}

func createNugetResourceFromData(data *schemas.NugetFeedTypeResourceModel) (*feeds.NuGetFeed, error) {
	feed, err := feeds.NewNuGetFeed(data.Name.ValueString(), data.FeedUri.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.DownloadAttempts = int(data.DownloadAttempts.ValueInt64())
	feed.DownloadRetryBackoffSeconds = int(data.DownloadRetryBackoffSeconds.ValueInt64())
	feed.EnhancedMode = data.IsEnhancedMode.ValueBool()

	var packageAcquisitionLocationOptions []string
	for _, element := range data.PackageAcquisitionLocationOptions.Elements() {
		packageAcquisitionLocationOptions = append(packageAcquisitionLocationOptions, element.(types.String).ValueString())
	}

	feed.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptions
	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()
	feed.Username = data.Username.ValueString()

	return feed, nil
}

func updateDataFromNugetFeed(data *schemas.NugetFeedTypeResourceModel, spaceId string, feed *feeds.NuGetFeed) {
	data.DownloadAttempts = types.Int64Value(int64(feed.DownloadAttempts))
	data.DownloadRetryBackoffSeconds = types.Int64Value(int64(feed.DownloadRetryBackoffSeconds))
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}
	data.IsEnhancedMode = types.BoolValue(feed.EnhancedMode)

	packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
	for i, option := range feed.PackageAcquisitionLocationOptions {
		packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
	}

	var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
	data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	data.ID = types.StringValue(feed.GetID())
}

func (*nugetFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
