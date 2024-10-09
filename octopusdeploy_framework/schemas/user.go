package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	UserResourceDescription = "user"
)

type UserSchema struct{}

var _ EntitySchema = UserSchema{}

func UserObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                     types.StringType,
		"username":               types.StringType,
		"can_password_be_edited": types.BoolType,
		"display_name":           types.StringType,
		"email_address":          types.StringType,
		"is_active":              types.BoolType,
		"is_requestor":           types.BoolType,
		"is_service":             types.BoolType,
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
			"space_id": GetUserSpaceIdDatasourceSchema(),
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
		"can_password_be_edited": GetBooleanDatasourceAttribute("Specifies whether or not the password can be edited.", true),
		"display_name":           GetDisplayNameDatasourceSchema(),
		"email_address":          GetEmailAddressDatasourceSchema(),
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
						Computed:    true,
					},
					"claim": datasourceSchema.SetNestedAttribute{
						Description: "The claim associated with the identity.",
						Computed:    true,
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

func GetUserSpaceIdDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description:        "The space ID associated with this user.",
		Optional:           true,
		DeprecationMessage: "This attribute is deprecated and will be removed in a future release. Users are not scoped to spaces, meaning providing a space ID will not affect the result.",
	}
}

func GetFilterDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A filter search by username, display name or email",
		Optional:    true,
	}
}

func GetDisplayNameDatasourceSchema() datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The display name of this resource.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Required: true,
	}

	return s
}

func GetEmailAddressDatasourceSchema() datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The email address of this resource.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Optional: true,
	}

	return s
}

func IdentityObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"provider": types.StringType,
		"claim":    types.SetType{ElemType: types.ObjectType{AttrTypes: IdentityClaimObjectType()}},
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

func MapFromUser(u *users.User) UserTypeResourceModel {
	var user UserTypeResourceModel
	user.ID = types.StringValue(u.ID)
	user.Username = types.StringValue(u.Username)
	user.CanPasswordBeEdited = types.BoolValue(u.CanPasswordBeEdited)
	user.DisplayName = types.StringValue(u.DisplayName)
	user.EmailAddress = types.StringValue(u.EmailAddress)
	user.IsActive = types.BoolValue(u.IsActive)
	user.IsRequestor = types.BoolValue(u.IsRequestor)
	user.IsService = types.BoolValue(u.IsService)
	user.Identity = types.SetValueMust(types.ObjectType{AttrTypes: IdentityObjectType()}, MapIdentities(u.Identities))

	return user
}

type UserTypeResourceModel struct {
	Username            types.String `tfsdk:"username"`
	CanPasswordBeEdited types.Bool   `tfsdk:"can_password_be_edited"`
	DisplayName         types.String `tfsdk:"display_name"`
	EmailAddress        types.String `tfsdk:"email_address"`
	IsActive            types.Bool   `tfsdk:"is_active"`
	IsRequestor         types.Bool   `tfsdk:"is_requestor"`
	IsService           types.Bool   `tfsdk:"is_service"`
	Identity            types.Set    `tfsdk:"identity"`

	ResourceModel
}
