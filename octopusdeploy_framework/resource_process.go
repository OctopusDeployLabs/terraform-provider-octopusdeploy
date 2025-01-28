package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &processResource{}

type processResource struct {
	*Config
}

func NewProcessResource() resource.Resource {
	return &processResource{}
}

func (r *processResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessResourceName)
}

func (r *processResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessSchema{}.GetResourceSchema()
}

func (r *processResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	ownerId := data.OwnerID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("creating process for owner: %s", ownerId))

	client := r.Config.Client
	// Empty process is created as part of project creation
	project, err := projects.GetByID(client, spaceId, ownerId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process, unable to find associated project", err.Error())
		return
	}

	process, err := deployments.GetDeploymentProcessByID(client, spaceId, project.DeploymentProcessID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process", err.Error())
		return
	}

	data.ID = types.StringValue(process.ID)
	data.SpaceID = types.StringValue(process.SpaceID)

	tflog.Info(ctx, fmt.Sprintf("process created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process (%s)", data.ID))

	client := r.Config.Client
	spaceId := data.SpaceID.ValueString()
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to read process", err.Error())
		return
	}

	data.ID = types.StringValue(process.ID)
	data.SpaceID = types.StringValue(process.SpaceID)

	tflog.Info(ctx, fmt.Sprintf("process read (%s)", process.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating process (%s)", data.ID))

	client := r.Config.Client
	spaceId := data.SpaceID.ValueString()
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	process.ProjectID = data.OwnerID.ValueString()

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process", err.Error())
		return
	}

	data.ID = types.StringValue(updatedProcess.ID)
	data.SpaceID = types.StringValue(updatedProcess.SpaceID)

	tflog.Info(ctx, fmt.Sprintf("process updated (%s)", updatedProcess.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("deleting process (%s)", data.ID))

	client := r.Config.Client
	spaceId := data.SpaceID.ValueString()
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	process.ProjectID = data.OwnerID.ValueString()
	process.SpaceID = data.SpaceID.ValueString()
	process.Steps = []*deployments.DeploymentStep{}

	_, err = deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
