package octopusdeploy_framework

import (
	"context"
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
	return
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
	return nil
}

func flattenDeploymentTarget(ctx context.Context, deploymentTarget machines.DeploymentTarget, model schemas.DeploymentTargetModel) schemas.DeploymentTargetModel {
	return model
}
