package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type azureWebAppDeploymentTargetResource struct {
	*Config
}

func NewAzureWebAppDeploymentTargetResource() resource.Resource {
	return &azureWebAppDeploymentTargetResource{}
}

var _ resource.ResourceWithImportState = &azureWebAppDeploymentTargetResource{}

func (r *azureWebAppDeploymentTargetResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("azure_web_app_deployment_target")
}

func (r *azureWebAppDeploymentTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AzureWebAppDeploymentTargetSchema{}.GetResourceSchema()
}

func (r *azureWebAppDeploymentTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *azureWebAppDeploymentTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.AzureWebAppDeploymentTargetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Azure Web App deployment target", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	deploymentTarget := expandAzureWebAppDeploymentTarget(ctx, plan)
	createdDeploymentTarget, err := machines.Add(r.Config.Client, deploymentTarget)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Amazon Web App deployment target", err.Error())
		return
	}

	state := flattenAzureWebAppDeploymentTarget(ctx, *createdDeploymentTarget, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *azureWebAppDeploymentTargetResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	return
}

func (r *azureWebAppDeploymentTargetResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	return
}

func (r *azureWebAppDeploymentTargetResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	return
}

func (*azureWebAppDeploymentTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandAzureWebAppDeploymentTarget(ctx context.Context, model schemas.AzureWebAppDeploymentTargetModel) *machines.DeploymentTarget {
	deploymentTarget := expandDeploymentTarget(ctx, model.DeploymentTargetModel)

	endpoint := machines.NewAzureWebAppEndpoint()
	endpoint.AccountID = model.AccountId.ValueString()
	endpoint.ResourceGroupName = model.ResourceGroupName.ValueString()
	endpoint.WebAppName = model.WebAppName.ValueString()
	endpoint.WebAppSlotName = model.WebAppSlotName.ValueString()
	deploymentTarget.Endpoint = endpoint

	return deploymentTarget
}

func flattenAzureWebAppDeploymentTarget(ctx context.Context, deploymentTarget machines.DeploymentTarget, model schemas.AzureWebAppDeploymentTargetModel) schemas.AzureWebAppDeploymentTargetModel {
	flattenedDeploymentTarget := flattenDeploymentTarget(ctx, deploymentTarget, model.DeploymentTargetModel)

	model.ID = flattenedDeploymentTarget.ID
	model.Endpoint = flattenedDeploymentTarget.Endpoint
	model.Environments = flattenedDeploymentTarget.Environments
	model.HasLatestCalamari = flattenedDeploymentTarget.HasLatestCalamari
	model.HealthStatus = flattenedDeploymentTarget.HealthStatus
	model.IsDisabled = flattenedDeploymentTarget.IsDisabled
	model.IsInProcess = flattenedDeploymentTarget.IsInProcess
	model.MachinePolicyId = flattenedDeploymentTarget.MachinePolicyId
	model.Name = flattenedDeploymentTarget.Name
	model.OperatingSystem = flattenedDeploymentTarget.OperatingSystem
	model.Roles = flattenedDeploymentTarget.Roles
	model.ShellName = flattenedDeploymentTarget.ShellName
	model.ShellVersion = flattenedDeploymentTarget.ShellVersion
	model.SpaceId = flattenedDeploymentTarget.SpaceId
	model.Status = flattenedDeploymentTarget.Status
	model.StatusSummary = flattenedDeploymentTarget.StatusSummary
	model.TenantedDeploymentParticipation = flattenedDeploymentTarget.TenantedDeploymentParticipation
	model.Tenants = flattenedDeploymentTarget.Tenants
	model.TenantTags = flattenedDeploymentTarget.TenantTags
	model.Thumbprint = flattenedDeploymentTarget.Thumbprint
	model.Uri = flattenedDeploymentTarget.Uri

	endpointResource, _ := machines.ToEndpointResource(deploymentTarget.Endpoint)

	model.AccountId = types.StringValue(endpointResource.AccountID)
	model.ResourceGroupName = types.StringValue(endpointResource.ResourceGroupName)
	model.WebAppName = types.StringValue(endpointResource.WebAppName)
	model.WebAppSlotName = types.StringValue(endpointResource.WebAppSlotName)

	return model
}
