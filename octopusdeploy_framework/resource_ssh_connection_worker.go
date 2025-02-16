package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type sshConnectionWorkerResource struct {
	*Config
}

func NewSSHConnectionWorkerResource() resource.Resource {
	return &sshConnectionWorkerResource{}
}

var _ resource.ResourceWithImportState = &sshConnectionWorkerResource{}

func (r *sshConnectionWorkerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("ssh_connection_worker")
}

func (r *sshConnectionWorkerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.SSHConnectionWorkerSchema{}.GetResourceSchema()
}

func (r *sshConnectionWorkerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *sshConnectionWorkerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.SSHConnectionWorkerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	worker := createSSHConnectionWorkerResource(ctx, data)

	tflog.Info(ctx, fmt.Sprintf("creating SSH connection worker: %s", data.Name.ValueString()))

	client := r.Config.Client
	createdWorker, err := workers.Add(client, worker)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create SSH connection worker", err.Error())
		return
	}

	updateDataFromSSHConnectionWorker(ctx, data, data.SpaceID.ValueString(), createdWorker)

	tflog.Info(ctx, fmt.Sprintf("SSH connection worker created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshConnectionWorkerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.SSHConnectionWorkerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading SSH connection worker (%s)", data.ID))

	client := r.Config.Client
	worker, err := workers.GetByID(client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "SSH connection worker"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load SSH connection worker", err.Error())
		}
		return
	}

	if worker.Endpoint.GetCommunicationStyle() != "Ssh" {
		resp.Diagnostics.AddError("unable to load SSH connection worker", "found resource is not SSH connection worker")
		return
	}

	updateDataFromSSHConnectionWorker(ctx, data, data.SpaceID.ValueString(), worker)

	tflog.Info(ctx, fmt.Sprintf("SSH connection worker read (%s)", worker.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshConnectionWorkerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.SSHConnectionWorkerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating SSH connection worker '%s'", data.ID.ValueString()))

	worker := createSSHConnectionWorkerResource(ctx, data)
	worker.ID = state.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("updating SSH connection worker (%s)", data.ID))

	client := r.Config.Client
	updatedWorker, err := workers.Update(client, worker)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update SSH connection worker", err.Error())
		return
	}

	updateDataFromSSHConnectionWorker(ctx, data, state.SpaceID.ValueString(), updatedWorker)

	tflog.Info(ctx, fmt.Sprintf("SSH connection worker updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshConnectionWorkerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.SSHConnectionWorkerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := workers.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete SSH connection worker", err.Error())
		return
	}
}

func createSSHConnectionWorkerResource(ctx context.Context, data *schemas.SSHConnectionWorkerResourceModel) *machines.Worker {
	endpoint := machines.NewSSHEndpoint(data.Host.ValueString(), int(data.Port.ValueInt64()), data.Fingerprint.ValueString())
	endpoint.AccountID = data.AccountId.ValueString()
	endpoint.DotNetCorePlatform = data.DotnetPlatform.ValueString()
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

func updateDataFromSSHConnectionWorker(ctx context.Context, data *schemas.SSHConnectionWorkerResourceModel, spaceId string, worker *machines.Worker) {
	data.ID = types.StringValue(worker.ID)
	data.SpaceID = types.StringValue(spaceId)
	data.Name = types.StringValue(worker.Name)
	data.IsDisabled = types.BoolValue(worker.IsDisabled)
	data.MachinePolicyID = types.StringValue(worker.MachinePolicyID)
	data.WorkerPoolIDs, _ = types.SetValueFrom(ctx, types.StringType, worker.WorkerPoolIDs)

	endpoint := worker.Endpoint.(*machines.SSHEndpoint)
	data.AccountId = types.StringValue(endpoint.AccountID)
	data.Host = types.StringValue(endpoint.Host)
	data.Port = types.Int64Value(int64(endpoint.Port))
	data.Fingerprint = types.StringValue(endpoint.Fingerprint)
	if endpoint.ProxyID != "" {
		data.ProxyID = types.StringValue(endpoint.ProxyID)
	}
}

func (*sshConnectionWorkerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
