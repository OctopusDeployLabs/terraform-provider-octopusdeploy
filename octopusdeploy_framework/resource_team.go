package octopusdeploy_framework

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	var plan schemas.TeamTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team := teams.NewTeam(plan.Name.ValueString())
	team.CanBeDeleted = plan.CanBeDeleted.ValueBool()
	team.CanBeRenamed = plan.CanBeRenamed.ValueBool()
	team.CanChangeMembers = plan.CanChangeMembers.ValueBool()
	team.CanChangeRoles = plan.CanChangeRoles.ValueBool()
	team.Description = plan.Description.ValueString()
	team.SpaceID = plan.SpaceId.ValueString()
	team.Description = plan.Description.ValueString()

	var userIds []string
	plan.Users.ElementsAs(ctx, &userIds, false)

	team.MemberUserIDs = userIds                                                        // TODO: Verify this is correct
	team.ExternalSecurityGroups = mapExternalSecurityGroup(plan.ExternalSecurityGroups) // TODO: Verify this is correct

	newTeam, err := r.Config.Client.Teams.Add(team)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create team", err.Error())
		return
	}

	scopedUserRoles, err := createScopedUserRoles(ctx, r.Client, &plan, newTeam)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create roles", err.Error())
		return
	}

	updatePlan(newTeam, scopedUserRoles, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func updatePlan(newTeam *teams.Team, scopedUserRoles []*userroles.ScopedUserRole, plan *schemas.TeamTypeResourceModel) {
	plan.ID = types.StringValue(newTeam.ID)
	plan.Name = types.StringValue(newTeam.Name)
	plan.SpaceId = types.StringValue(newTeam.SpaceID)
	plan.CanBeDeleted = types.BoolValue(newTeam.CanBeDeleted)
	plan.CanBeRenamed = types.BoolValue(newTeam.CanBeRenamed)
	plan.CanChangeMembers = types.BoolValue(newTeam.CanChangeMembers)
	plan.CanChangeRoles = types.BoolValue(newTeam.CanChangeRoles)
	plan.ExternalSecurityGroups = schemas.MapToExternalSecurityGroupsDatasourceModel(newTeam.ExternalSecurityGroups)
	plan.Users = basetypes.SetValue(util.FlattenStringList(newTeam.MemberUserIDs))
	plan.UserRole = MapToScopedUserRoleResourceModel(scopedUserRoles)
}

func createScopedUserRoles(ctx context.Context, client *client.Client, plan *schemas.TeamTypeResourceModel, team *teams.Team) ([]*userroles.ScopedUserRole, error) {
	newScopedUserRoles := make([]*userroles.ScopedUserRole, 0, len(plan.UserRole))
	for _, planUserRole := range plan.UserRole {
		scopedUserRole := userroles.NewScopedUserRole(planUserRole.UserRoleID.ValueString())
		scopedUserRole.TeamID = team.ID
		scopedUserRole.SpaceID = planUserRole.SpaceID.ValueString()

		// Verify this is correct
		scopedUserRole.ID = planUserRole.ID.ValueString()
		scopedUserRole.EnvironmentIDs, _ = util.SetToStringArray(ctx, planUserRole.EnvironmentIDs) // TODO: Handle diagnostics
		scopedUserRole.ProjectGroupIDs, _ = util.SetToStringArray(ctx, planUserRole.ProjectGroupIDs)
		scopedUserRole.ProjectIDs, _ = util.SetToStringArray(ctx, planUserRole.ProjectIDs)
		scopedUserRole.TenantIDs, _ = util.SetToStringArray(ctx, planUserRole.TenantIDs)

		newScopedUserRoles = append(newScopedUserRoles, scopedUserRole)
	}

	scopedUserRoles := make([]*userroles.ScopedUserRole, 0, len(newScopedUserRoles))
	for _, userRole := range newScopedUserRoles {
		createdScopedUserRole, err := client.ScopedUserRoles.Add(userRole)
		if err != nil {
			return []*userroles.ScopedUserRole{}, fmt.Errorf("error creating user role for team %s: %s", team.ID, err)
		}

		scopedUserRole, err := client.ScopedUserRoles.GetByID(createdScopedUserRole.ID)
		if err != nil {
			return []*userroles.ScopedUserRole{}, fmt.Errorf("error getting user role for team %s: %s", team.ID, err)
		}
		scopedUserRoles = append(scopedUserRoles, scopedUserRole)
	}

	return scopedUserRoles, nil
}

func MapToScopedUserRoleResourceModel(scopedUserRoles []*userroles.ScopedUserRole) []schemas.ScopedUserRoleResourceModel {
	models := make([]schemas.ScopedUserRoleResourceModel, len(scopedUserRoles))

	for i, scopedUserRole := range scopedUserRoles {
		models[i].UserRoleID = types.StringValue(scopedUserRole.UserRoleID)
		models[i].EnvironmentIDs = types.SetValueMust(types.StringType, util.ToValueSlice(scopedUserRole.EnvironmentIDs))
		models[i].ID = types.StringValue(scopedUserRole.ID)
		models[i].ProjectGroupIDs = types.SetValueMust(types.StringType, util.ToValueSlice(scopedUserRole.ProjectGroupIDs))
		models[i].ProjectIDs = types.SetValueMust(types.StringType, util.ToValueSlice(scopedUserRole.ProjectIDs))
		models[i].TenantIDs = types.SetValueMust(types.StringType, util.ToValueSlice(scopedUserRole.TenantIDs))
		models[i].SpaceID = types.StringValue(scopedUserRole.SpaceID)
		models[i].TeamID = types.StringValue(scopedUserRole.TeamID)
	}

	return models
}

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

	updateTeam(&data, team) // Move userrole mapping to mapToTeamResourceModel
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
