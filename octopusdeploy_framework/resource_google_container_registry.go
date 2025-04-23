package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

	googleContainerRegistryFeed, err := createContainerRegistryFeedResourceFromGoogleData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Google Container Registry feed: %s", googleContainerRegistryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, googleContainerRegistryFeed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create Google Container Registry feed", err.Error())
		return
	}

	updateGoogleDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.GoogleContainerRegistry))

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
			resp.Diagnostics.AddError("unable to load Google Container Registry feed", err.Error())
		}
		return
	}

	googleContainerRegistry := feed.(*feeds.GoogleContainerRegistry)
	updateGoogleDataFromDockerContainerRegistryFeed(data, data.SpaceID.ValueString(), googleContainerRegistry)

	tflog.Info(ctx, fmt.Sprintf("Google Container Registry feed read (%s)", googleContainerRegistry.GetID()))
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

	feed, err := createContainerRegistryFeedResourceFromGoogleData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load Google Container Registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Google Container Registry feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update Google Container Registry feed", err.Error())
		return
	}

	updateGoogleDataFromDockerContainerRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.GoogleContainerRegistry))

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
		resp.Diagnostics.AddError("unable to delete Google Container Registry feed", err.Error())
		return
	}
}

func createContainerRegistryFeedResourceFromGoogleData(data *schemas.GoogleContainerRegistryFeedTypeResourceModel) (*feeds.GoogleContainerRegistry, error) {
	var oidc *feeds.GoogleContainerRegistryOidcAuthentication

	if data.OidcAuthentication != nil {
		oidc = &feeds.GoogleContainerRegistryOidcAuthentication{
			Audience:    data.OidcAuthentication.Audience.ValueString(),
			SubjectKeys: util.ExpandStringList(data.OidcAuthentication.SubjectKey),
		}
	}

	feed, err := feeds.NewGoogleContainerRegistry(data.Name.ValueString(), data.Username.ValueString(), core.NewSensitiveValue(data.Password.ValueString()), oidc)

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

func updateGoogleDataFromDockerContainerRegistryFeed(data *schemas.GoogleContainerRegistryFeedTypeResourceModel, spaceId string, feed *feeds.GoogleContainerRegistry) {
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

	if feed.OidcAuthentication != nil {
		data.OidcAuthentication = &schemas.GoogleContainerRegistryOidcAuthenticationResourceModel{
			Audience:   types.StringValue(feed.OidcAuthentication.Audience),
			SubjectKey: util.FlattenStringList(feed.OidcAuthentication.SubjectKeys),
		}
	}
}

func (*googleContainerRegistryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
