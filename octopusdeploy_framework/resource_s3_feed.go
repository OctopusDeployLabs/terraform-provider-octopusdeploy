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

type s3FeedTypeResource struct {
	*Config
}

func NewS3FeedResource() resource.Resource {
	return &s3FeedTypeResource{}
}

var _ resource.ResourceWithImportState = &s3FeedTypeResource{}

func (r *s3FeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("s3_feed")
}

func (r *s3FeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.S3FeedSchema{}.GetResourceSchema()
}

func (r *s3FeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *s3FeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.S3FeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	feed, err := createS3ResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating S3 feed: %s", feed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to create S3 feed", err.Error())
		return
	}

	updateDataFromS3Feed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.S3Feed))

	tflog.Info(ctx, fmt.Sprintf("S3 feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *s3FeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.S3FeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading S3 feed (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "S3 feed"); err != nil {
			resp.Diagnostics.AddError("unable to load S3 feed", err.Error())
		}
		return
	}

	loadedFeed := feed.(*feeds.S3Feed)
	updateDataFromS3Feed(data, data.SpaceID.ValueString(), loadedFeed)

	tflog.Info(ctx, fmt.Sprintf("S3 feed read (%s)", loadedFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *s3FeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.S3FeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating S3 feed '%s'", data.ID.ValueString()))

	feed, err := createS3ResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		resp.Diagnostics.AddError("unable to load S3 feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating S3 feed (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		resp.Diagnostics.AddError("unable to update S3 feed", err.Error())
		return
	}

	updateDataFromS3Feed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.S3Feed))

	tflog.Info(ctx, fmt.Sprintf("S3 feed updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *s3FeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.S3FeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete S3 feed", err.Error())
		return
	}
}

func createS3ResourceFromData(data *schemas.S3FeedTypeResourceModel) (*feeds.S3Feed, error) {
	feed, err := feeds.NewS3Feed(data.Name.ValueString(), data.AccessKey.ValueString(), core.NewSensitiveValue(data.SecretKey.ValueString()), data.UseMachineCredentials.ValueBool())
	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()

	feed.Username = data.Username.ValueString()
	feed.Password = core.NewSensitiveValue(data.Password.ValueString())
	feed.SpaceID = data.SpaceID.ValueString()

	return feed, nil
}

func updateDataFromS3Feed(data *schemas.S3FeedTypeResourceModel, spaceId string, feed *feeds.S3Feed) {
	data.UseMachineCredentials = types.BoolValue(feed.UseMachineCredentials)
	if feed.AccessKey != "" {
		data.AccessKey = types.StringValue(feed.AccessKey)
	}
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	if feed.Username != "" {
		data.Username = types.StringValue(feed.Username)
	}

	data.ID = types.StringValue(feed.ID)
}

func (*s3FeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
