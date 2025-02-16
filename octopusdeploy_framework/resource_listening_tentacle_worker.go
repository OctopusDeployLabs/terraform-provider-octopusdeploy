package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"net/url"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type listeningTentacleWorkerResource struct {
	*Config
}

func NewListeningTentacleWorkerResource() resource.Resource {
	return &listeningTentacleWorkerResource{}
}

var _ resource.ResourceWithImportState = &listeningTentacleWorkerResource{}

func (r *listeningTentacleWorkerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("listening_tentacle_worker")
}

func (r *listeningTentacleWorkerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ListeningTentacleWorkerSchema{}.GetResourceSchema()
}

func (r *listeningTentacleWorkerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *listeningTentacleWorkerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ListeningTentacleWorkerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	worker := createListeningTentacleWorkerResource(ctx, data)

	tflog.Info(ctx, fmt.Sprintf("creating listening tentacle worker: %s", data.Name.ValueString()))

	client := r.Config.Client
	createdWorker, err := workers.Add(client, worker)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create listening tentacle worker", err.Error())
		return
	}

	updateDataFromListeningTentacleWorker(ctx, data, data.SpaceID.ValueString(), createdWorker)

	tflog.Info(ctx, fmt.Sprintf("listening tentacle worker created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *listeningTentacleWorkerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ListeningTentacleWorkerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading listening tentacle worker (%s)", data.ID))

	client := r.Config.Client
	worker, err := workers.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "listening tentacle worker"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load listening tentacle worker", err.Error())
		}
		return
	}

	if worker.Endpoint.GetCommunicationStyle() != "TentaclePassive" {
		resp.Diagnostics.AddError("unable to load listening tentacle worker", "found resource is not listening tentacle worker")
		return
	}

	updateDataFromListeningTentacleWorker(ctx, data, data.SpaceID.ValueString(), worker)

	tflog.Info(ctx, fmt.Sprintf("listening tentacle worker read (%s)", worker.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *listeningTentacleWorkerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.ListeningTentacleWorkerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating listening tentacle worker '%s'", data.ID.ValueString()))

	worker := createListeningTentacleWorkerResource(ctx, data)
	worker.ID = state.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("updating listening tentacle worker (%s)", data.ID))

	client := r.Config.Client
	updatedWorker, err := workers.Update(client, worker)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update listening tentacle worker", err.Error())
		return
	}

	updateDataFromListeningTentacleWorker(ctx, data, state.SpaceID.ValueString(), updatedWorker)

	tflog.Info(ctx, fmt.Sprintf("listening tentacle worker updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *listeningTentacleWorkerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.ListeningTentacleWorkerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := workers.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete listening tentacle worker", err.Error())
		return
	}
}

func createListeningTentacleWorkerResource(ctx context.Context, data *schemas.ListeningTentacleWorkerResourceModel) *machines.Worker {
	uri, _ := url.Parse(data.Uri.ValueString())
	endpoint := machines.NewListeningTentacleEndpoint(uri, data.Thumbprint.ValueString())
	endpoint.ProxyID = data.ProxyID.ValueString()

	worker := machines.NewWorker(data.Name.ValueString(), endpoint)
	worker.SpaceID = data.SpaceID.ValueString()
	worker.IsDisabled = data.IsDisabled.ValueBool()
	worker.MachinePolicyID = data.MachinePolicyID.ValueString()

	if !data.WorkerPoolIDs.IsNull() {
		var workerPools []string
		data.WorkerPoolIDs.ElementsAs(ctx, &workerPools, false)
		worker.WorkerPoolIDs = workerPools
	}

	return worker
}

func updateDataFromListeningTentacleWorker(ctx context.Context, data *schemas.ListeningTentacleWorkerResourceModel, spaceId string, worker *machines.Worker) {
	data.ID = types.StringValue(worker.ID)
	data.SpaceID = types.StringValue(spaceId)
	data.Name = types.StringValue(worker.Name)
	data.IsDisabled = types.BoolValue(worker.IsDisabled)
	data.MachinePolicyID = types.StringValue(worker.MachinePolicyID)
	data.WorkerPoolIDs, _ = types.SetValueFrom(ctx, types.StringType, worker.WorkerPoolIDs)

	endpoint := worker.Endpoint.(*machines.ListeningTentacleEndpoint)
	data.Uri = types.StringValue(endpoint.URI.String())
	data.Thumbprint = types.StringValue(endpoint.Thumbprint)
	if endpoint.ProxyID != "" {
		data.ProxyID = types.StringValue(endpoint.ProxyID)
	}
}

func (*listeningTentacleWorkerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
