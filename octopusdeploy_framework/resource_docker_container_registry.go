package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dockerContainerRegistryFeedTypeResource struct {
	*Config
}

func NewDockerContainerRegistryFeedResource() resource.Resource {
	return &dockerContainerRegistryFeedTypeResource{}
}

func (r *dockerContainerRegistryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("docker_container_registry")
}

func (r *dockerContainerRegistryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:  schemas.GetDockerContainerRegistryFeedResourceSchema(),
		Description: "This resource manages a Docker Container Registry in Octopus Deploy.",
	}
}

func (r *dockerContainerRegistryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *dockerContainerRegistryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.DockerContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dockerContainerRegistryFeed, err := createDockerContainerRegistryFeedResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Docker Container Registry feed: %s", dockerContainerRegistryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, dockerContainerRegistryFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create docker container registry feed", err.Error())
		return
	}

	updateDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Docker Container Registry feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerContainerRegistryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.DockerContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Docker Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load docker container registry feed", err.Error())
		return
	}

	dockerContainerRegistry := feed.(*feeds.DockerContainerRegistry)
	updateDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), dockerContainerRegistry)

	tflog.Info(ctx, fmt.Sprintf("Docker Container Registry feed read (%s)", dockerContainerRegistry.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerContainerRegistryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.DockerContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating docker container registry feed '%s'", data.ID.ValueString()))

	feed, err := createDockerContainerRegistryFeedResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load docker container registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Docker Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update docker container registry feed", err.Error())
		return
	}

	updateDataFromDockerContainerRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.DockerContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Docker Container Registry feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerContainerRegistryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.DockerContainerRegistryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete docker container registry feed", err.Error())
		return
	}
}

func createDockerContainerRegistryFeedResourceFromData(data *schemas.DockerContainerRegistryFeedTypeResourceModel) (*feeds.DockerContainerRegistry, error) {
	feed, err := feeds.NewDockerContainerRegistry(data.Name.ValueString())
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
	feed.APIVersion = data.APIVersion.ValueString()
	feed.RegistryPath = data.RegistryPath.ValueString()

	return feed, nil
}

func updateDataFromDockerContainerRegistryFeed(data *schemas.DockerContainerRegistryFeedTypeResourceModel, spaceId string, feed *feeds.DockerContainerRegistry) {
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

	packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
	for i, option := range feed.PackageAcquisitionLocationOptions {
		packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
	}

	var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
	data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	data.ID = types.StringValue(feed.ID)
}
