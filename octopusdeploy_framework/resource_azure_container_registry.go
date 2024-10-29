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

type azureContainerRegistryFeedTypeResource struct {
	*Config
}

func NewAzureContainerRegistryFeedResource() resource.Resource {
	return &azureContainerRegistryFeedTypeResource{}
}

var _ resource.ResourceWithImportState = &azureContainerRegistryFeedTypeResource{}

func (r *azureContainerRegistryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("azure_container_registry")
}

func (r *azureContainerRegistryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GoogleContainerRegistryFeedSchema{}.GetResourceSchema()
}

func (r *azureContainerRegistryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *azureContainerRegistryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.AzureContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dockerContainerRegistryFeed, err := createDockerContainerRegistryFeedResourceFromAzureData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Azure Container Registry feed: %s", dockerContainerRegistryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, dockerContainerRegistryFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create Azure Container Registry feed", err.Error())
		return
	}

	updateAzureDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Azure Container Registry feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureContainerRegistryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.AzureContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Azure Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "azure container registry feed"); err != nil {
			resp.Diagnostics.AddError("unable to load Azure Container Registry feed", err.Error())
		}
		return
	}

	dockerContainerRegistry := feed.(*feeds.DockerContainerRegistry)
	updateAzureDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), dockerContainerRegistry)

	tflog.Info(ctx, fmt.Sprintf("Azure Container Registry feed read (%s)", dockerContainerRegistry.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureContainerRegistryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.AzureContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating Azure Container Registry feed '%s'", data.ID.ValueString()))

	feed, err := createDockerContainerRegistryFeedResourceFromAzureData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load Azure Container Registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Azure Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update Azure Container Registry feed", err.Error())
		return
	}

	updateAzureDataFromDockerContainerRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Google Container Registry feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *azureContainerRegistryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.AzureContainerRegistryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete Azure Container Registry feed", err.Error())
		return
	}
}

func createDockerContainerRegistryFeedResourceFromAzureData(data *schemas.AzureContainerRegistryFeedTypeResourceModel) (*feeds.DockerContainerRegistry, error) {
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

func updateAzureDataFromDockerContainerRegistryFeed(data *schemas.AzureContainerRegistryFeedTypeResourceModel, spaceId string, feed *feeds.DockerContainerRegistry) {
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

func (*azureContainerRegistryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
