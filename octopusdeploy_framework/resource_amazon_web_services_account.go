package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type amazonWebServicesAccountResource struct {
	*Config
}

func NewAmazonWebServicesAccountResource() resource.Resource {
	return &amazonWebServicesAccountResource{}
}

var _ resource.ResourceWithImportState = &amazonWebServicesAccountResource{}

func (r *amazonWebServicesAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("aws_account")
}

func (r *amazonWebServicesAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AmazonWebServicesAccountSchema{}.GetResourceSchema()
}

func (r *amazonWebServicesAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *amazonWebServicesAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.AmazonWebServicesAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Amazon Web Services account", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	account := expandAmazonWebServicesAccount(ctx, plan)
	createdAccount, err := accounts.Add(r.Config.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Amazon Web Services account", err.Error())
		return
	}

	state := flattenAmazonWebServicesAccount(ctx, createdAccount.(*accounts.AmazonWebServicesAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *amazonWebServicesAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *amazonWebServicesAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *amazonWebServicesAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (*amazonWebServicesAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandAmazonWebServicesAccount(ctx context.Context, model schemas.AmazonWebServicesAccountModel) *accounts.AmazonWebServicesAccount {
	var accountName = model.Name.ValueString()
	var accountAccessKey = model.AccessKey.ValueString()
	var accountSecretKey = core.NewSensitiveValue(model.SecretKey.ValueString())

	account, _ := accounts.NewAmazonWebServicesAccount(accountName, accountAccessKey, accountSecretKey)

	account.SetID(model.ID.ValueString())
	account.SetDescription(model.Description.ValueString())
	account.SetEnvironmentIDs(expandStringList(model.Environments))
	account.SetSpaceID(model.SpaceId.ValueString())
	account.SetTenantedDeploymentMode(core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString()))
	account.SetTenantIDs(expandStringList(model.Tenants))
	account.SetTenantTags(expandStringList(model.TenantTags))

	return account
}
func flattenAmazonWebServicesAccount(ctx context.Context, account *accounts.AmazonWebServicesAccount, model schemas.AmazonWebServicesAccountModel) schemas.AmazonWebServicesAccountModel {
	model.ID = types.StringValue(account.GetID())
	model.AccessKey = types.StringValue(account.AccessKey)
	model.Description = types.StringValue(account.GetDescription())
	model.Environments = flattenStringList(account.GetEnvironmentIDs(), model.Environments)
	model.Name = types.StringValue(account.GetName())
	model.SpaceId = types.StringValue(account.GetSpaceID())
	model.TenantedDeploymentParticipation = types.StringValue(string(account.GetTenantedDeploymentMode()))
	model.Tenants = flattenStringList(account.GetTenantIDs(), model.Tenants)
	model.TenantTags = flattenStringList(account.GetTenantTags(), model.TenantTags)

	// Note: We don't flatten the secret key as it's sensitive and not returned by the API

	return model
}
