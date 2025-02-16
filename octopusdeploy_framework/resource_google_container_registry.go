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

type googleContainerRegistryFeedTypeResource struct {
	*Config
}

func NewGoogleContainerRegistryFeedResource() resource.Resource {
	return &googleContainerRegistryFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &googleContainerRegistryFeedTypeResource{}

func (r *googleContainerRegistryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("google_container_registry")
}

func (r *googleContainerRegistryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GoogleContainerRegistryFeedSchema{}.GetResourceSchema()
}

func (r *googleContainerRegistryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *googleContainerRegistryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.GoogleContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dockerContainerRegistryFeed, err := createDockerContainerRegistryFeedResourceFromGoogleData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Google Container Registry feed: %s", dockerContainerRegistryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, dockerContainerRegistryFeed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create Google Container Registry feed", err.Error())
		return
	}

	updateGoogleDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Google Container Registry feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *googleContainerRegistryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.GoogleContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Google Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "google container registry feed"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load Google Container Registry feed", err.Error())
		}
		return
	}

	dockerContainerRegistry := feed.(*feeds.DockerContainerRegistry)
	updateGoogleDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), dockerContainerRegistry)

	tflog.Info(ctx, fmt.Sprintf("Google Container Registry feed read (%s)", dockerContainerRegistry.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *googleContainerRegistryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.GoogleContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating Google Container Registry feed '%s'", data.ID.ValueString()))

	feed, err := createDockerContainerRegistryFeedResourceFromGoogleData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load Google Container Registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Google Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update Google Container Registry feed", err.Error())
		return
	}

	updateGoogleDataFromDockerContainerRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Google Container Registry feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *googleContainerRegistryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.GoogleContainerRegistryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete Google Container Registry feed", err.Error())
		return
	}
}

func createDockerContainerRegistryFeedResourceFromGoogleData(data *schemas.GoogleContainerRegistryFeedTypeResourceModel) (*feeds.DockerContainerRegistry, error) {
	feed, err := feeds.NewDockerContainerRegistry(data.Name.ValueString())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()
	feed.FeedURI = data.FeedUri.ValueString()
	feed.PackageAcquisitionLocationOptions = nil
	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()
	feed.Username = data.Username.ValueString()
	feed.APIVersion = data.APIVersion.ValueString()
	feed.RegistryPath = data.RegistryPath.ValueString()

	return feed, nil
}

func updateGoogleDataFromDockerContainerRegistryFeed(data *schemas.GoogleContainerRegistryFeedTypeResourceModel, spaceId string, feed *feeds.DockerContainerRegistry) {
	data.FeedUri = types.StringValue(feed.FeedURI)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.APIVersion != "" {
		data.APIVersion = types.StringValue(feed.APIVersion)
	}
	if feed.RegistryPath != "" {
		data.RegistryPath = types.StringValue(feed.RegistryPath)
	}
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}

	data.ID = types.StringValue(feed.ID)
}

func (*googleContainerRegistryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
