package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &genericOidcAccountResource{}
var _ resource.ResourceWithImportState = &genericOidcAccountResource{}

type genericOidcAccountResource struct {
	*Config
}

func NewGenericOidcResource() resource.Resource {
	return &genericOidcAccountResource{}
}

func (r *genericOidcAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("generic_oidc_account")
}

func (r *genericOidcAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GenericOidcAccountSchema{}.GetResourceSchema()
}

func (r *genericOidcAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}
func (r *genericOidcAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.GenericOidcAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating generic oidc account", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	account := expandGenericOidcAccountResource(ctx, plan)
	createdAccount, err := accounts.Add(r.Client, account)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "Error creating generic oidc account", err.Error())
		return
	}

	state := flattenGenericOidcAccountResource(ctx, createdAccount.(*accounts.GenericOIDCAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *genericOidcAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.GenericOidcAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := accounts.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "genericOidcAccountResource"); err != nil {
			util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "unable to load generic oidc account", err.Error())
		}
		return
	}

	newState := flattenGenericOidcAccountResource(ctx, account.(*accounts.GenericOIDCAccount), state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *genericOidcAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.GenericOidcAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account := expandGenericOidcAccountResource(ctx, plan)
	updatedAccount, err := accounts.Update(r.Client, account)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "Error updating generic oidc account", err.Error())
		return
	}

	state := flattenGenericOidcAccountResource(ctx, updatedAccount.(*accounts.GenericOIDCAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *genericOidcAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.GenericOidcAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := accounts.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, r.Config.SystemInfo, "Error deleting generic oidc account", err.Error())
		return
	}
}

func (r *genericOidcAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	accountID := req.ID

	account, err := accounts.GetByID(r.Client, r.Client.GetSpaceID(), accountID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading generic oidc account",
			fmt.Sprintf("Unable to read generic oidc account with ID %s: %s", accountID, err.Error()),
		)
		return
	}

	genericOidcAccount, ok := account.(*accounts.GenericOIDCAccount)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected account type",
			fmt.Sprintf("Expected generic oidc account, got: %T", account),
		)
		return
	}

	state := schemas.GenericOidcAccountResourceModel{
		SpaceID:                         types.StringValue(genericOidcAccount.GetSpaceID()),
		Name:                            types.StringValue(genericOidcAccount.GetName()),
		Description:                     types.StringValue(genericOidcAccount.GetDescription()),
		TenantedDeploymentParticipation: types.StringValue(string(genericOidcAccount.GetTenantedDeploymentMode())),
		Environments:                    flattenStringList(genericOidcAccount.GetEnvironmentIDs(), types.ListNull(types.StringType)),
		Tenants:                         flattenStringList(genericOidcAccount.GetTenantIDs(), types.ListNull(types.StringType)),
		TenantTags:                      flattenStringList(genericOidcAccount.TenantTags, types.ListNull(types.StringType)),
		ExecutionSubjectKeys:            flattenStringList(genericOidcAccount.DeploymentSubjectKeys, types.ListNull(types.StringType)),
		Audience:                        types.StringValue(genericOidcAccount.Audience),
	}
	state.ID = types.StringValue(genericOidcAccount.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func expandGenericOidcAccountResource(ctx context.Context, model schemas.GenericOidcAccountResourceModel) *accounts.GenericOIDCAccount {
	account, _ := accounts.NewGenericOIDCAccount(model.Name.ValueString())

	account.SetID(model.ID.ValueString())
	account.SetDescription(model.Description.ValueString())
	account.SetSpaceID(model.SpaceID.ValueString())
	account.SetEnvironmentIDs(util.ExpandStringList(model.Environments))
	account.SetTenantedDeploymentMode(core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString()))
	account.SetTenantIDs(util.ExpandStringList(model.Tenants))
	account.SetTenantTags(util.ExpandStringList(model.TenantTags))
	account.DeploymentSubjectKeys = util.ExpandStringList(model.ExecutionSubjectKeys)
	account.Audience = model.Audience.ValueString()

	return account
}

func flattenGenericOidcAccountResource(ctx context.Context, account *accounts.GenericOIDCAccount, model schemas.GenericOidcAccountResourceModel) schemas.GenericOidcAccountResourceModel {
	model.ID = types.StringValue(account.GetID())
	model.SpaceID = types.StringValue(account.GetSpaceID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.TenantedDeploymentParticipation = types.StringValue(string(account.GetTenantedDeploymentMode()))

	model.Environments = util.FlattenStringList(account.GetEnvironmentIDs())
	model.Tenants = util.FlattenStringList(account.GetTenantIDs())
	model.TenantTags = util.FlattenStringList(account.TenantTags)

	model.ExecutionSubjectKeys = util.FlattenStringList(account.DeploymentSubjectKeys)
	model.Audience = types.StringValue(account.Audience)

	return model
}
