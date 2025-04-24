package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	var plan schemas.AzureSubscriptionAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Azure subscription account", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	account := expandAzureSubscriptionAccount(ctx, plan)
	createdAccount, err := accounts.Add(r.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Azure subscription account", err.Error())
		return
	}

	state := flattenAzureSubscriptionAccount(ctx, createdAccount.(*accounts.AzureSubscriptionAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *azureSubscriptionAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.AzureSubscriptionAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := accounts.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "azureSubscriptionAccountResource"); err != nil {
			resp.Diagnostics.AddError("unable to load Azure subscription account", err.Error())
		}
		return
	}

	newState := flattenAzureSubscriptionAccount(ctx, account.(*accounts.AzureSubscriptionAccount), state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)

}

func (r *azureSubscriptionAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.AzureSubscriptionAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account := expandAzureSubscriptionAccount(ctx, plan)
	updatedAccount, err := accounts.Update(r.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Azure subscription account", err.Error())
		return
	}

	state := flattenAzureSubscriptionAccount(ctx, updatedAccount.(*accounts.AzureSubscriptionAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *azureSubscriptionAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.AzureSubscriptionAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := accounts.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Azure subscription account", err.Error())
		return
	}

}

func (*azureSubscriptionAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandAzureSubscriptionAccount(ctx context.Context, model schemas.AzureSubscriptionAccountModel) *accounts.AzureSubscriptionAccount {
	var name = model.Name.ValueString()
	var subscriptionId, _ = uuid.Parse(model.SubscriptionID.ValueString())

	account, _ := accounts.NewAzureSubscriptionAccount(name, subscriptionId)

	account.SetID(model.ID.ValueString())
	account.AzureEnvironment = model.AzureEnvironment.ValueString()
	account.CertificateBytes = core.NewSensitiveValue(model.Certificate.ValueString())
	account.CertificateThumbprint = model.CertificateThumbprint.ValueString()
	account.SetDescription(model.Description.ValueString())
	account.SetEnvironmentIDs(expandStringList(model.Environments))
	account.ManagementEndpoint = model.ManagementEndpoint.ValueString()
	account.SetName(model.Name.ValueString())
	account.SetSpaceID(model.SpaceID.ValueString())
	account.StorageEndpointSuffix = model.StorageEndpointSuffix.ValueString()
	account.SetTenantedDeploymentMode(core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString()))
	account.SetTenantTags(expandStringList(model.TenantTags))
	account.SetTenantIDs(expandStringList(model.Tenants))

	return account

}

func flattenAzureSubscriptionAccount(ctx context.Context, account *accounts.AzureSubscriptionAccount, model schemas.AzureSubscriptionAccountModel) schemas.AzureSubscriptionAccountModel {
	model.ID = types.StringValue(account.GetID())
	model.AzureEnvironment = types.StringValue(account.AzureEnvironment)
	model.CertificateThumbprint = types.StringValue(account.CertificateThumbprint)
	model.Description = types.StringValue(account.GetDescription())
	model.Environments = flattenStringList(account.GetEnvironmentIDs(), model.Environments)
	model.ManagementEndpoint = types.StringValue(account.ManagementEndpoint)
	model.Name = types.StringValue(account.GetName())
	model.SubscriptionID = types.StringValue(account.SubscriptionID.String())
	model.SpaceID = types.StringValue(account.GetSpaceID())
	model.StorageEndpointSuffix = types.StringValue(account.StorageEndpointSuffix)
	model.TenantedDeploymentParticipation = types.StringValue(string(account.GetTenantedDeploymentMode()))
	model.Tenants = flattenStringList(account.GetTenantIDs(), model.Tenants)
	model.TenantTags = flattenStringList(account.TenantTags, model.TenantTags)

	// Note: We don't flatten the certificate as it's sensitive and not returned by the API

	return model
}
