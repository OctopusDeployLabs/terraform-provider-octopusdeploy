package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &usernamePasswordAccountResource{}
var _ resource.ResourceWithImportState = &usernamePasswordAccountResource{}

type usernamePasswordAccountResource struct {
	*Config
}

func NewUsernamePasswordAccountResource() resource.Resource {
	return &usernamePasswordAccountResource{}
}

func (r *usernamePasswordAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("username_password_account")
}

func (r *usernamePasswordAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.UsernamePasswordAccountSchema{}.GetResourceSchema()
}

func (r *usernamePasswordAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}
func (r *usernamePasswordAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.UsernamePasswordAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating username password account", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	account := expandUsernamePasswordAccount(ctx, plan)
	createdAccount, err := accounts.Add(r.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating username password account", err.Error())
		return
	}

	state := flattenUsernamePasswordAccount(ctx, createdAccount.(*accounts.UsernamePasswordAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *usernamePasswordAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.UsernamePasswordAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := accounts.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "usernamePasswordAccountResource"); err != nil {
			resp.Diagnostics.AddError("unable to load username password account", err.Error())
		}
		return
	}

	newState := flattenUsernamePasswordAccount(ctx, account.(*accounts.UsernamePasswordAccount), state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *usernamePasswordAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.UsernamePasswordAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	account := expandUsernamePasswordAccount(ctx, plan)
	updatedAccount, err := accounts.Update(r.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating username password account", err.Error())
		return
	}

	state := flattenUsernamePasswordAccount(ctx, updatedAccount.(*accounts.UsernamePasswordAccount), plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *usernamePasswordAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.UsernamePasswordAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := accounts.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting username password account", err.Error())
		return
	}
}

func (r *usernamePasswordAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	accountID := req.ID

	account, err := accounts.GetByID(r.Client, r.Client.GetSpaceID(), accountID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading username password account",
			fmt.Sprintf("Unable to read username password account with ID %s: %s", accountID, err.Error()),
		)
		return
	}

	usernamePasswordAccount, ok := account.(*accounts.UsernamePasswordAccount)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected account type",
			fmt.Sprintf("Expected username password account, got: %T", account),
		)
		return
	}

	state := schemas.UsernamePasswordAccountResourceModel{
		SpaceID:                         types.StringValue(usernamePasswordAccount.GetSpaceID()),
		Name:                            types.StringValue(usernamePasswordAccount.GetName()),
		Description:                     types.StringValue(usernamePasswordAccount.GetDescription()),
		Username:                        types.StringValue(usernamePasswordAccount.GetUsername()),
		TenantedDeploymentParticipation: types.StringValue(string(usernamePasswordAccount.GetTenantedDeploymentMode())),
		Environments:                    flattenStringList(usernamePasswordAccount.GetEnvironmentIDs(), types.ListNull(types.StringType)),
		Tenants:                         flattenStringList(usernamePasswordAccount.GetTenantIDs(), types.ListNull(types.StringType)),
		TenantTags:                      flattenStringList(usernamePasswordAccount.TenantTags, types.ListNull(types.StringType)),
		Password:                        types.StringNull(),
	}
	state.ID = types.StringValue(usernamePasswordAccount.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func expandUsernamePasswordAccount(ctx context.Context, model schemas.UsernamePasswordAccountResourceModel) *accounts.UsernamePasswordAccount {
	account, _ := accounts.NewUsernamePasswordAccount(model.Name.ValueString())

	account.SetID(model.ID.ValueString())
	account.SetDescription(model.Description.ValueString())
	account.SetSpaceID(model.SpaceID.ValueString())
	account.SetUsername(model.Username.ValueString())
	account.SetPassword(core.NewSensitiveValue(model.Password.ValueString()))
	account.SetEnvironmentIDs(expandStringList(model.Environments))
	account.SetTenantedDeploymentMode(core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString()))
	account.SetTenantIDs(expandStringList(model.Tenants))
	account.SetTenantTags(expandStringList(model.TenantTags))

	return account
}

func flattenUsernamePasswordAccount(ctx context.Context, account *accounts.UsernamePasswordAccount, model schemas.UsernamePasswordAccountResourceModel) schemas.UsernamePasswordAccountResourceModel {
	model.ID = types.StringValue(account.GetID())
	model.SpaceID = types.StringValue(account.GetSpaceID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.Username = types.StringValue(account.GetUsername())
	model.TenantedDeploymentParticipation = types.StringValue(string(account.GetTenantedDeploymentMode()))

	model.Environments = flattenStringList(account.GetEnvironmentIDs(), model.Environments)
	model.Tenants = flattenStringList(account.GetTenantIDs(), model.Tenants)
	model.TenantTags = flattenStringList(account.TenantTags, model.TenantTags)

	// Note: We don't flatten the password as it's sensitive and not returned by the API

	return model
}

func expandStringList(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var result []string
	list.ElementsAs(context.Background(), &result, false)
	if len(result) == 0 {
		return nil
	}

	return result
}

func flattenStringList(slice []string, currentList types.List) types.List {
	if len(slice) == 0 && currentList.IsNull() {
		return types.ListNull(types.StringType)
	}
	if slice == nil {
		return types.ListNull(types.StringType)
	}

	valueSlice := make([]attr.Value, len(slice))
	for i, s := range slice {
		valueSlice[i] = types.StringValue(s)
	}

	return types.ListValueMust(types.StringType, valueSlice)
}

func expandStringSet(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var result []string
	set.ElementsAs(context.Background(), &result, false)
	if len(result) == 0 {
		return nil
	}

	return result
}

func flattenStringSet(slice []string, currentSet types.Set) types.Set {
	if len(slice) == 0 && currentSet.IsNull() {
		return types.SetNull(types.StringType)
	}
	if slice == nil {
		return types.SetNull(types.StringType)
	}

	valueSlice := make([]attr.Value, len(slice))
	for i, s := range slice {
		valueSlice[i] = types.StringValue(s)
	}

	return types.SetValueMust(types.StringType, valueSlice)
}
