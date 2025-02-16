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

type awsElasticContainerRegistryFeedTypeResource struct {
	*Config
}

var _ resource.ResourceWithImportState = &awsElasticContainerRegistryFeedTypeResource{}

func NewAwsElasticContainerRegistryFeedResource() resource.Resource {
	return &awsElasticContainerRegistryFeedTypeResource{}
}

func (r *awsElasticContainerRegistryFeedTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("aws_elastic_container_registry")
}

func (r *awsElasticContainerRegistryFeedTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AwsElasticContainerRegistrySchema{}.GetResourceSchema()
}

func (r *awsElasticContainerRegistryFeedTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *awsElasticContainerRegistryFeedTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.AwsElasticContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	awsElasticContainerRegistryFeed, err := createAwsElasticContainerRegistryResourceFromData(data)
	if err != nil {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Aws Elastic Container Registry: %s", awsElasticContainerRegistryFeed.GetName()))

	client := r.Config.Client
	createdFeed, err := feeds.Add(client, awsElasticContainerRegistryFeed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create aws elastic container registry", err.Error())
		return
	}

	updateDataFromAwsElasticContainerRegistryFeed(data, data.SpaceID.ValueString(), createdFeed.(*feeds.AwsElasticContainerRegistry))

	data.ID = types.StringValue(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("Aws Elastic Container Registry feed created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *awsElasticContainerRegistryFeedTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.AwsElasticContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Aws Elastic Container Registry (%s)", data.ID))

	client := r.Config.Client
	feed, err := feeds.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "aws elastic container registry"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load aws elastic container registry", err.Error())
		}
		return
	}

	awsElasticContainerRegistryFeed := feed.(*feeds.AwsElasticContainerRegistry)
	updateDataFromAwsElasticContainerRegistryFeed(data, data.SpaceID.ValueString(), awsElasticContainerRegistryFeed)

	tflog.Info(ctx, fmt.Sprintf("Aws Elastic Container Registry read (%s)", awsElasticContainerRegistryFeed.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *awsElasticContainerRegistryFeedTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.AwsElasticContainerRegistryFeedTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating aws elastic container registry feed '%s'", data.ID.ValueString()))

	feed, err := createAwsElasticContainerRegistryResourceFromData(data)
	feed.ID = state.ID.ValueString()
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load aws elastic container registry feed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating Aws Elastic Container Registry (%s)", data.ID))

	client := r.Config.Client
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update aws elastic container registry feed", err.Error())
		return
	}

	updateDataFromAwsElasticContainerRegistryFeed(data, state.SpaceID.ValueString(), updatedFeed.(*feeds.AwsElasticContainerRegistry))

	tflog.Info(ctx, fmt.Sprintf("Aws Elastic Container Registry updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *awsElasticContainerRegistryFeedTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.AwsElasticContainerRegistryFeedTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := feeds.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete aws elastic container registry feed", err.Error())
		return
	}
}

func createAwsElasticContainerRegistryResourceFromData(data *schemas.AwsElasticContainerRegistryFeedTypeResourceModel) (*feeds.AwsElasticContainerRegistry, error) {
	feed, err := feeds.NewAwsElasticContainerRegistry(data.Name.ValueString(), data.AccessKey.ValueString(), core.NewSensitiveValue(data.SecretKey.ValueString()), data.Region.ValueString())

	if err != nil {
		return nil, err
	}

	feed.ID = data.ID.ValueString()

	var packageAcquisitionLocationOptions []string
	for _, element := range data.PackageAcquisitionLocationOptions.Elements() {
		packageAcquisitionLocationOptions = append(packageAcquisitionLocationOptions, element.(types.String).ValueString())
	}

	feed.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptions
	feed.SpaceID = data.SpaceID.ValueString()

	return feed, nil
}

func updateDataFromAwsElasticContainerRegistryFeed(data *schemas.AwsElasticContainerRegistryFeedTypeResourceModel, spaceId string, feed *feeds.AwsElasticContainerRegistry) {
	data.AccessKey = types.StringValue(feed.AccessKey)
	data.Name = types.StringValue(feed.Name)
	data.SpaceID = types.StringValue(spaceId)
	data.Region = types.StringValue(feed.Region)

	packageAcquisitionLocationOptionsList := make([]attr.Value, len(feed.PackageAcquisitionLocationOptions))
	for i, option := range feed.PackageAcquisitionLocationOptions {
		packageAcquisitionLocationOptionsList[i] = types.StringValue(option)
	}

	var packageAcquisitionLocationOptionsListValue, _ = types.ListValue(types.StringType, packageAcquisitionLocationOptionsList)
	data.PackageAcquisitionLocationOptions = packageAcquisitionLocationOptionsListValue
	data.ID = types.StringValue(feed.GetID())
}

func (*awsElasticContainerRegistryFeedTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
