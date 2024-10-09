package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithImportState = &userTypeResource{}

type userTypeResource struct {
	*Config
}

func NewUserResource() resource.Resource { return &environmentTypeResource{} }

func (r *userTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("user")
}

func (r *userTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.UserSchema{}.GetResourceSchema()
}

func (r *userTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *userTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.UserTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newUser := users.NewUser(data.Username.ValueString(), data.DisplayName.ValueString())
	newUser.Password = data.Password.ValueString()
	newUser.CanPasswordBeEdited = data.CanPasswordBeEdited.ValueBool()
	newUser.EmailAddress = data.EmailAddress.ValueString()
	newUser.IsActive = data.IsActive.ValueBool()
	newUser.IsRequestor = data.IsRequestor.ValueBool()
	newUser.IsService = data.IsService.ValueBool()
	if len(data.Identity.Elements()) > 0 {
		newUser.Identities = mapIdentities(data.Identity)
	}

	user, err := users.Add(r.Config.Client, newUser)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create user", err.Error())
		return
	}

	updateUser(&data, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *userTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *userTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func updateUser(data *schemas.UserTypeResourceModel, user *users.User) {
	data.ID = types.StringValue(user.ID)
	data.Username = types.StringValue(user.Username)
	data.Password = types.StringValue(user.Password)
	data.CanPasswordBeEdited = types.BoolValue(user.CanPasswordBeEdited)
	data.DisplayName = types.StringValue(user.DisplayName)
	data.EmailAddress = types.StringValue(user.EmailAddress)
	data.IsActive = types.BoolValue(user.IsActive)
	data.IsRequestor = types.BoolValue(user.IsRequestor)
	data.IsService = types.BoolValue(user.IsService)
	data.Identity = types.SetValueMust(types.ObjectType{AttrTypes: schemas.IdentityObjectType()}, schemas.MapIdentities(user.Identities))
}

func mapIdentities(identities types.Set) []users.Identity {
	result := make([]users.Identity, 0, len(identities.Elements()))
	for _, identityElem := range identities.Elements() {
		identityObj := identityElem.(types.Object)
		identityAttrs := identityObj.Attributes()

		identity := users.Identity{}
		if v, ok := identityAttrs["provider"].(types.String); ok && !v.IsNull() {
			identity.IdentityProviderName = v.ValueString()
		}

		if v, ok := identityAttrs["claim"].(types.List); ok && !v.IsNull() {
			identity.Claims = mapIdentityClaims(v)
		}
	}

	return result
}

func mapIdentityClaims(identityClaims types.List) map[string]users.IdentityClaim {
	result := map[string]users.IdentityClaim{}
	for _, identityClaimElem := range identityClaims.Elements() {
		identityClaimObj := identityClaimElem.(types.Object)
		identityClaimAttrs := identityClaimObj.Attributes()

		identityClaim := users.IdentityClaim{}
		var name string
		// TODO: Handle error if name isn't found?
		if v, ok := identityClaimAttrs["name"].(types.String); ok && !v.IsNull() {
			name = v.ValueString()
		}

		if v, ok := identityClaimAttrs["is_identifying_claim"].(types.Bool); ok && !v.IsNull() {
			identityClaim.IsIdentifyingClaim = v.ValueBool()
		}

		if v, ok := identityClaimAttrs["value"].(types.String); ok && !v.IsNull() {
			identityClaim.Value = v.ValueString()
		}

		result[name] = identityClaim
	}

	return result
}
