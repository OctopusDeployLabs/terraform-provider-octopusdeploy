package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type deploymentTargetResource struct {
	*Config
}

func NewDeploymentTargetResource() resource.Resource {
	return &deploymentTargetResource{}
}

var _ resource.ResourceWithImportState = &deploymentTargetResource{}

func (r *deploymentTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("deployment_target")
}

func (r *deploymentTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AzureWebAppDeploymentTargetSchema{}.GetResourceSchema()
}

func (r *deploymentTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *deploymentTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *deploymentTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *deploymentTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *deploymentTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (*deploymentTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandDeploymentTarget(ctx context.Context, model schemas.DeploymentTargetModel) *machines.DeploymentTarget {
	name := model.Name.ValueString()
	endpoint := expandAzureWebAppDeploymentEndpoint(ctx, model)
	environments := expandStringList(model.Environments)
	roles := expandStringList(model.Roles)

	deploymentTarget := machines.NewDeploymentTarget(name, endpoint, environments, roles)

	deploymentTarget.ID = model.ID.ValueString()
	deploymentTarget.HasLatestCalamari = model.HasLatestCalamari.ValueBool()
	deploymentTarget.HealthStatus = model.HealthStatus.ValueString()
	deploymentTarget.IsDisabled = model.IsDisabled.ValueBool()
	deploymentTarget.IsInProcess = model.IsInProcess.ValueBool()
	deploymentTarget.MachinePolicyID = model.MachinePolicyId.ValueString()
	deploymentTarget.OperatingSystem = model.OperatingSystem.ValueString()
	deploymentTarget.ShellName = model.ShellName.ValueString()
	deploymentTarget.ShellVersion = model.ShellVersion.ValueString()
	deploymentTarget.SpaceID = model.SpaceId.ValueString()
	deploymentTarget.Status = model.Status.ValueString()
	deploymentTarget.StatusSummary = model.StatusSummary.ValueString()
	deploymentTarget.TenantedDeploymentMode = core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString())
	deploymentTarget.TenantIDs = expandStringList(model.Tenants)
	deploymentTarget.TenantTags = expandStringList(model.TenantTags)
	deploymentTarget.Thumbprint = model.Thumbprint.ValueString()
	deploymentTarget.URI = model.Uri.ValueString()

	return deploymentTarget
}

func expandAzureWebAppDeploymentEndpoint(_ context.Context, model schemas.DeploymentTargetModel) *machines.AzureWebAppEndpoint {
	endpoint := machines.NewAzureWebAppEndpoint()

	endpointAttributes := model.Endpoint.Attributes()

	endpoint.AccountID = endpointAttributes["account_id"].String()
	endpoint.ID = endpointAttributes["id"].String()
	endpoint.ResourceGroupName = endpointAttributes["resource_group_name"].String()
	endpoint.WebAppName = endpointAttributes["web_app_name"].String()
	endpoint.WebAppSlotName = endpointAttributes["web_app_slot_name"].String()

	return endpoint
}

func flattenDeploymentTarget(ctx context.Context, deploymentTarget machines.DeploymentTarget, model schemas.DeploymentTargetModel) schemas.DeploymentTargetModel {
	return model
}
