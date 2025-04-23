package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type azureSubscriptionAccountResource struct {
	*Config
}

func NewAzureSubscriptionAccountResource() resource.Resource {
	return &azureSubscriptionAccountResource{}
}

var _ resource.ResourceWithImportState = &azureSubscriptionAccountResource{}

func (r *azureSubscriptionAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("azure_subscription_account")
}

func (r *azureSubscriptionAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AzureSubscriptionAccountSchema{}.GetResourceSchema()
}

func (r *azureSubscriptionAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *azureSubscriptionAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	return
}

func (r *azureSubscriptionAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *azureSubscriptionAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *azureSubscriptionAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (*azureSubscriptionAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandAzureSubscriptionAccount(ctx context.Context, model schemas.AzureSubscriptionAccountModel) *accounts.AzureSubscriptionAccount {
	return nil
}

func flattenAmazonWebServicesAccount(ctx context.Context, account *accounts.AmazonWebServicesAccount, model schemas.AzureSubscriptionAccountModel) schemas.AzureSubscriptionAccountModel {
	return model
}
