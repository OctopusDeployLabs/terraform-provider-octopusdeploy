package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

var _ resource.ResourceWithImportState = &teamTypeResource{}

type teamTypeResource struct {
	*Config
}

func NewTeamResource() resource.Resource { return &teamTypeResource{} }

func (r *teamTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("team")
}

func (r *teamTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TeamSchema{}.GetResourceSchema()
}

func (r *teamTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *teamTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *teamTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.TeamTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newTeam := teams.NewTeam(data.Name.ValueString())
	newTeam.CanBeDeleted = data.CanBeDeleted.ValueBool()
	newTeam.CanBeRenamed = data.CanBeRenamed.ValueBool()
	newTeam.CanChangeMembers = data.CanChangeMembers.ValueBool()
	newTeam.CanChangeRoles = data.CanChangeRoles.ValueBool()
	newTeam.Description = data.Description.ValueString()
	newTeam.SpaceID = data.SpaceId.ValueString()
	newTeam.Description = data.Description.ValueString()

	var userIds []string
	data.Users.ElementsAs(ctx, &userIds, false)

	newTeam.MemberUserIDs = userIds                                                        // TODO: Verify this is correct
	newTeam.ExternalSecurityGroups = mapExternalSecurityGroup(data.ExternalSecurityGroups) // TODO: Verify this is correct

	team, err := r.Config.Client.Teams.Add(newTeam)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create team", err.Error())
		return
	}

	// Octopus doesn't allow creating inactive users. To mimic creating an inactive user, we need to update the newly created user.
	//if !data.IsActive.ValueBool() {
	//	user.IsActive = data.IsActive.ValueBool()
	//	user, err = users.Update(r.Config.Client, user)
	//}

	err = resourceTeamUpdateUserRoles(r.Client, ctx, data, team)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create roles", err.Error())
		return
	}

	roles := data.UserRole
	data = schemas.MapToTeamsResourceModel(team)

	data.UserRole = roles
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func resourceTeamUpdateUserRoles(client *client.Client, ctx context.Context, d schemas.TeamTypeResourceModel, team *teams.Team) error {
	log.Printf("[INFO] updating team user roles (%s)", d.ID)
	//if d.UserRole.Equal(team.UserRoles) {
	//	log.Printf("[INFO] user role has changes (%s)", d.ID)
	//	o, n := d.GetChange("user_role")
	//	if o == nil {
	//		o = new(schema.Set)
	//	}
	//	if n == nil {
	//		n = new(schema.Set)
	//	}
	//
	//	os := o.(*schema.Set)
	//	ns := n.(*schema.Set)
	//	remove := expandUserRoles(team, os.Difference(ns).List())
	//	add := expandUserRoles(team, ns.Difference(os).List())
	//
	//	if len(remove) > 0 || len(add) > 0 {
	//		log.Printf("[INFO] user role found diff (%s)", d.ID)
	//		//client := m.(*client.Client)
	//		if len(remove) > 0 {
	//			log.Printf("[INFO] removing user roles from team (%s)", d.ID)
	//			for _, userRole := range remove {
	//				if userRole.ID != "" {
	//					err := client.ScopedUserRoles.DeleteByID(userRole.ID)
	//					if err != nil {
	//						apiError := err.(*core.APIError)
	//						if apiError.StatusCode != 404 {
	//							// It's already been deleted, maybe mixing with the independent resource?
	//							return fmt.Errorf("error removing user role %s from team %s: %s", userRole.ID, team.ID, err)
	//						}
	//					}
	//				}
	//			}
	//		}
	//		if len(add) > 0 {
	//			log.Printf("[INFO] adding new user roles to team (%s)", d.ID)
	//			for _, userRole := range add {
	//				_, err := client.ScopedUserRoles.Add(userRole)
	//				if err != nil {
	//					return fmt.Errorf("error creating user role for team %s: %s", team.ID, err)
	//				}
	//			}
	//		}
	//	}
	//}
	return nil
}

func expandUserRoles(team *teams.Team, userRoles []interface{}) []*userroles.ScopedUserRole {
	values := make([]*userroles.ScopedUserRole, 0, len(userRoles))
	for _, rawUserRole := range userRoles {
		userRole := rawUserRole.(map[string]interface{})
		scopedUserRole := userroles.NewScopedUserRole(userRole["user_role_id"].(string))
		scopedUserRole.TeamID = team.ID
		scopedUserRole.SpaceID = userRole["space_id"].(string)

		if v, ok := userRole["id"]; ok {
			scopedUserRole.ID = v.(string)
		} else {
			scopedUserRole.ID = ""
		}

		if v, ok := userRole["environment_ids"]; ok {
			scopedUserRole.EnvironmentIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_group_ids"]; ok {
			scopedUserRole.ProjectGroupIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_ids"]; ok {
			scopedUserRole.ProjectIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["tenant_ids"]; ok {
			scopedUserRole.TenantIDs = getSliceFromTerraformTypeList(v)
		}
		values = append(values, scopedUserRole)
	}
	return values
}

func getSliceFromTerraformTypeList(list interface{}) []string {
	if list == nil {
		return nil
	}

	if v, ok := list.([]string); ok {
		return v
	}

	terraformList, ok := list.([]interface{})
	if !ok {
		terraformSet, ok := list.(*schema.Set)
		if ok {
			terraformList = terraformSet.List()
		} else {
			// It's not a list or set type
			return nil
		}
	}
	var newSlice []string
	for _, v := range terraformList {
		if v != nil {
			newSlice = append(newSlice, v.(string))
		}
	}
	return newSlice
}

//func mapIdentities(identities types.Set) []users.Identity {
//	result := make([]users.Identity, 0, len(identities.Elements()))
//	for _, identityElem := range identities.Elements() {
//		identityObj := identityElem.(types.Object)
//		identityAttrs := identityObj.Attributes()
//
//		identity := users.Identity{}
//		if v, ok := identityAttrs["provider"].(types.String); ok && !v.IsNull() {
//			identity.IdentityProviderName = v.ValueString()
//		}
//
//		if v, ok := identityAttrs["claim"].(types.Set); ok && !v.IsNull() {
//			identity.Claims = mapIdentityClaims(v)
//		}
//		result = append(result, identity)
//	}
//
//	return result
//}

func mapExternalSecurityGroup(externalSecurityGroups types.List) []core.NamedReferenceItem {
	expandedExternalSecurityGroups := make([]core.NamedReferenceItem, 0, len(externalSecurityGroups.Elements()))
	for _, externalSecurityGroupElem := range externalSecurityGroups.Elements() {
		groupObj := externalSecurityGroupElem.(types.Object)
		groupAttrs := groupObj.Attributes()

		group := core.NamedReferenceItem{}

		if v, ok := groupAttrs["display_id_and_name"].(types.Bool); ok && !v.IsNull() {
			group.DisplayIDAndName = v.ValueBool()
		}

		if v, ok := groupAttrs["display_name"].(types.String); ok && !v.IsNull() {
			group.DisplayName = v.ValueString()
		}

		if v, ok := groupAttrs["id"].(types.String); ok && !v.IsNull() {
			group.ID = v.ValueString()
		}

		expandedExternalSecurityGroups = append(expandedExternalSecurityGroups, group)
	}
	return expandedExternalSecurityGroups
}

func (r *teamTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.TeamTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.Client.Teams.GetByID(data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "team"); err != nil {
			resp.Diagnostics.AddError("unable to load team", err.Error())
		}
		return
	}

	updateTeam(&data, team)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateTeam(s *schemas.TeamTypeResourceModel, team *teams.Team) {

}

func (r *teamTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.UserTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := users.GetByID(r.Config.Client, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load user", err.Error())
		return
	}

	updatedUser := users.NewUser(data.Username.ValueString(), data.DisplayName.ValueString())
	updatedUser.ID = user.ID
	updatedUser.Password = data.Password.ValueString()
	updatedUser.EmailAddress = data.EmailAddress.ValueString()
	updatedUser.IsActive = data.IsActive.ValueBool()
	updatedUser.IsRequestor = data.IsRequestor.ValueBool()
	updatedUser.IsService = data.IsService.ValueBool()
	if len(data.Identity.Elements()) > 0 {
		updatedUser.Identities = mapIdentities(data.Identity)
	}

	updatedUser, err = users.Update(r.Config.Client, updatedUser)
	if err != nil {
		resp.Diagnostics.AddError("unable to update user", err.Error())
		return
	}

	updateUser(&data, updatedUser)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *teamTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.TeamTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.Teams.DeleteByID(data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete team", err.Error())
		return
	}
}
