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

type githubRepositoryFeedTypeResource struct {
	*Config
}

func NewGitHubRepositoryFeedResource() resource.Resource {
	return &githubRepositoryFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &githubRepositoryFeedTypeResource{}

func (r *githubRepositoryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("github_repository_feed")
}

func (r *githubRepositoryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GitHubRepositoryFeedSchema{}.GetResourceSchema()
}

func (r *githubRepositoryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *githubRepositoryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.GitHubRepositoryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	githubRepositoryFeed, err := createGitHubRepositoryResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating GitHub Repository feed: %s", githubRepositoryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, githubRepositoryFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create github repository feed", err.Error())
		return
	}

	updateDataFromGitHubRepositoryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.GitHubRepositoryFeed))

	data.ID = types.StringValue(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("GitHub Repository feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *githubRepositoryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.GitHubRepositoryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading GitHub Repository feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "github repository feed"); err != nil {
			resp.Diagnostics.AddError("unable to load github repository feed", err.Error())
		}
		return
	}

	githubRepositoryFeed := feed.(*feeds.GitHubRepositoryFeed)
	updateDataFromGitHubRepositoryFeed(data, data.SpaceID.ValueString(), githubRepositoryFeed)

	tflog.Info(ctx, fmt.Sprintf("GitHub Repository feed read (%s)", githubRepositoryFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *githubRepositoryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.GitHubRepositoryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating github repository feed '%s'", data.ID.ValueString()))

	feed, err := createGitHubRepositoryResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load github repository feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating GitHub Repository feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update github repository feed", err.Error())
		return
	}

	updateDataFromGitHubRepositoryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.GitHubRepositoryFeed))

	tflog.Info(ctx, fmt.Sprintf("GitHub Repository feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *githubRepositoryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.GitHubRepositoryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete github repository feed", err.Error())
		return
	}
}

func createGitHubRepositoryResourceFromData(data *schemas.GitHubRepositoryFeedTypeResourceModel) (*feeds.GitHubRepositoryFeed, error) {
	feed, err := feeds.NewGitHubRepositoryFeed(data.Name.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.DownloadAttempts = int(data.DownloadAttempts.ValueInt64())
	feed.DownloadRetryBackoffSeconds = int(data.DownloadRetryBackoffSeconds.ValueInt64())
	feed.FeedURI = data.FeedUri.ValueString()

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

func updateDataFromGitHubRepositoryFeed(data *schemas.GitHubRepositoryFeedTypeResourceModel, spaceId string, feed *feeds.GitHubRepositoryFeed) {
	data.DownloadAttempts = types.Int64Value(int64(feed.DownloadAttempts))
	data.DownloadRetryBackoffSeconds = types.Int64Value(int64(feed.DownloadRetryBackoffSeconds))
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}

	packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
	for i, option := range feed.PackageAcquisitionLocationOptions {
		packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
	}

	var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
	data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	data.ID = types.StringValue(feed.GetID())
}

func (*githubRepositoryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
