package schemas

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	UserResourceDescription = "user"
)

type UserSchema struct{}

var _ EntitySchema = UserSchema{}

func UserObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"username":      types.StringType,
		"password":      types.StringType,
		"display_name":  types.StringType,
		"email_address": types.StringType,
		"is_active":     types.BoolType,
		"is_requestor":  types.BoolType,
		"is_service":    types.BoolType,
		"identity": types.SetType{
			ElemType: types.ObjectType{AttrTypes: IdentityObjectType()},
		},
	}
}

func (u UserSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing users.",
		Attributes: map[string]datasourceSchema.Attribute{
			//request
			"ids":      GetQueryIDsDatasourceSchema(),
			"space_id": GetSpaceIdDatasourceSchema(UserResourceDescription, false),
			"filter":   GetFilterDatasourceSchema(),
			"skip":     GetQuerySkipDatasourceSchema(),
			"take":     GetQueryTakeDatasourceSchema(),

			//response
			"id": GetIdDatasourceSchema(true),
			"users": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: u.GetDatasourceSchemaAttributes(),
				},
			},
		},
	}
}

func (u UserSchema) GetDatasourceSchemaAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                     GetIdDatasourceSchema(true),
		"username":               GetUsernameDatasourceSchema(true),
		"password":               GetPasswordDatasourceSchema(false),
		"can_password_be_edited": GetBooleanDatasourceAttribute("Specifies whether or not the password can be edited.", true),
		"display_name":           GetDisplayNameDatasourceSchema(true),
		"email_address":          GetEmailAddressDatasourceSchema(false),
		"is_active":              GetBooleanDatasourceAttribute("Specifies whether or not the user is active.", true),
		"is_requestor":           GetBooleanDatasourceAttribute("Specifies whether or not the user is the requestor.", true),
		"is_service":             GetBooleanDatasourceAttribute("Specifies whether or not the user is a service account.", true),
		"identity": datasourceSchema.SetNestedAttribute{
			Description: "The identities associated with the user.",
			Optional:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					"provider": datasourceSchema.StringAttribute{
						Description: "The identity provider.",
					},
					"claim": datasourceSchema.SetNestedAttribute{
						Description: "The claim. // todo what is this",
						NestedObject: datasourceSchema.NestedAttributeObject{
							Attributes: map[string]datasourceSchema.Attribute{
								"name":                 GetNameDatasourceSchema(true),
								"is_identifying_claim": GetBooleanDatasourceAttribute("Specifies whether or not the claim is an identifying claim.", true),
								"value":                GetValueDatasourceSchema(true),
							},
						},
					},
				},
			},
		},
	}
}

func GetFilterDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A filter search by username, display name or email",
		Optional:    true,
	}
}

func IdentityObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"provider": types.StringType,
		"claim":    types.SetType{ElemType: types.ObjectType{AttrTypes: IdentityObjectType()}},
	}
}

func IdentityClaimObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                 types.StringType,
		"is_identifying_claim": types.BoolType,
		"value":                types.StringType,
	}
}

func (u UserSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

func MapIdentityClaims(claims map[string]users.IdentityClaim) []attr.Value {
	claimsList := make([]attr.Value, 0, len(claims))
	for key, claim := range claims {
		claimMap := map[string]attr.Value{
			"is_identifying_claim": types.BoolValue(claim.IsIdentifyingClaim),
			"name":                 types.StringValue(key),
			"value":                types.StringValue(claim.Value),
		}
		claimsList = append(claimsList, types.ObjectValueMust(IdentityClaimObjectType(), claimMap))
	}
	return claimsList
}

func MapIdentities(identities []users.Identity) []attr.Value {
	identitiesList := make([]attr.Value, 0, len(identities))
	for _, identity := range identities {
		identityMap := map[string]attr.Value{
			"provider": types.StringValue(identity.IdentityProviderName),
			"claim":    types.SetValueMust(types.ObjectType{AttrTypes: IdentityClaimObjectType()}, MapIdentityClaims(identity.Claims)),
		}
		identitiesList = append(identitiesList, types.ObjectValueMust(IdentityObjectType(), identityMap))
	}
	return identitiesList
}

func MapFromUser(ctx context.Context, u *users.User) UserTypeResourceModel {
	var user UserTypeResourceModel
	user.ID = types.StringValue(u.ID)
	user.Username = types.StringValue(u.Username)
	user.Password = types.StringValue(u.Password)
	user.DisplayName = types.StringValue(u.DisplayName)
	user.EmailAddress = types.StringValue(u.EmailAddress)
	user.IsActive = types.BoolValue(u.IsActive)
	user.IsRequestor = types.BoolValue(u.IsRequestor)
	user.IsService = types.BoolValue(u.IsService)
	user.Identity = types.SetValueMust(types.ObjectType{AttrTypes: IdentityObjectType()}, MapIdentities(u.Identities))
}

type UserTypeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	DisplayName  types.String `tfsdk:"display_name"`
	EmailAddress types.String `tfsdk:"email_address"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	IsRequestor  types.Bool   `tfsdk:"is_requestor"`
	IsService    types.Bool   `tfsdk:"is_service"`
	Identity     types.Set    `tfsdk:"identity"`

	ResourceModel
}
