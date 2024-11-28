package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ EntitySchema = TeamSchema{}

type TeamSchema struct{}

func (l TeamSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages lifecycles in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"can_be_deleted":          util.ResourceBool().Computed().Optional().Build(),
			"can_be_renamed":          util.ResourceBool().Computed().Optional().Build(),
			"can_change_members":      util.ResourceBool().Computed().Optional().Build(),
			"can_change_roles":        util.ResourceBool().Computed().Optional().Build(),
			"description":             util.ResourceString().Optional().Description("The user-friendly description of this team.").Build(),
			"external_security_group": getExternalSecurityGroupsAttributeResourceSchema(),
			"id":                      util.ResourceString().Computed().Optional().Description("The unique ID for this resource.").Build(),
			"name":                    util.ResourceString().Required().Description("The name of this team.").Build(),
			"space_id":                util.ResourceString().Computed().Optional().Description("The space associated with this team.").Build(),
			"users":                   util.ResourceSet(types.StringType).Computed().Optional().Description("A list of user IDs designated to be members of this team.").Build(),
		},
		Blocks: map[string]resourceSchema.Block{
			"user_role": resourceSchema.SetNestedBlock{
				Description: "The user roles associated with the team.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"id": util.ResourceString().Computed().Build(),
						//"id":                util.ResourceString().Optional().Computed().PlanModifiers(stringplanmodifier.UseStateForUnknown()).Description("The ID of the template parameter.").Build(),
						"environment_ids":   util.ResourceSet(types.StringType).Optional().Computed().PlanModifiers(setplanmodifier.UseStateForUnknown()).Build(),
						"project_group_ids": util.ResourceSet(types.StringType).Optional().Computed().PlanModifiers(setplanmodifier.UseStateForUnknown()).Build(),
						"project_ids":       util.ResourceSet(types.StringType).Optional().Computed().PlanModifiers(setplanmodifier.UseStateForUnknown()).Build(),
						"space_id":          util.ResourceString().Required().Build(),
						"team_id":           util.ResourceString().Computed().Build(),
						"tenant_ids":        util.ResourceSet(types.StringType).Optional().Computed().PlanModifiers(setplanmodifier.UseStateForUnknown()).Build(),
						"user_role_id":      util.ResourceString().Required().Build(),
					},
				},
			},
		},
	}
}

func (l TeamSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing teams.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":             util.DataSourceString().Computed().Description("An auto-generated identifier that includes the timestamp when this data source was last modified.").Build(),
			"ids":            util.DataSourceList(types.StringType).Optional().Description("A filter to search by a list of IDs..").Build(),
			"include_system": util.DataSourceBool().Optional().Description("A filter to include system teams.").Build(),
			"partial_name":   util.DataSourceString().Optional().Description("A filter to search by the partial match of a name.").Build(),
			"spaces":         util.DataSourceList(types.StringType).Optional().Description("A filter to search by a list of space IDs.").Build(),
			"skip":           util.DataSourceInt64().Optional().Description("A filter to specify the number of items to skip in the response.").Build(),
			"take":           util.DataSourceInt64().Optional().Description("A filter to specify the number of items to take (or return) in the response.").Build(),
			"teams":          getTeamsAttribute(),
		},
	}
}

func getTeamsAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed:    true,
		Description: "A list of teams that match the filter(s).",
		Optional:    false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"can_be_deleted":          util.DataSourceBool().Computed().Optional().Build(),
				"can_be_renamed":          util.DataSourceBool().Computed().Optional().Build(),
				"can_change_members":      util.DataSourceBool().Computed().Optional().Build(),
				"can_change_roles":        util.DataSourceBool().Computed().Optional().Build(),
				"description":             util.DataSourceString().Optional().Description("The user-friendly description of this team.").Build(),
				"external_security_group": getExternalSecurityGroupsAttribute(),
				"id":                      util.DataSourceString().Computed().Optional().Description("The unique ID for this resource.").Build(),
				"name":                    util.DataSourceString().Required().Description("The name of this team.").Build(),
				"space_id":                util.DataSourceString().Computed().Optional().Description("The space associated with this team.").Build(),
				"users":                   util.DataSourceSet(types.StringType).Computed().Optional().Description("A list of user IDs designated to be members of this team.").Build(),
			},
		},
	}
}

func getExternalSecurityGroupsAttributeResourceSchema() resourceSchema.ListNestedAttribute {
	return resourceSchema.ListNestedAttribute{
		Computed: false,
		Optional: true,
		NestedObject: resourceSchema.NestedAttributeObject{
			Attributes: map[string]resourceSchema.Attribute{
				"display_id_and_name": util.ResourceBool().Computed().Optional().Build(),
				"display_name":        util.ResourceString().Computed().Optional().Build(),
				"id":                  util.ResourceString().Computed().Optional().Description("The unique ID for this resource.").Build(),
			},
		},
	}
}

func getExternalSecurityGroupsAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: false,
		Optional: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"display_id_and_name": util.DataSourceBool().Computed().Optional().Build(),
				"display_name":        util.DataSourceString().Computed().Optional().Build(),
				"id":                  util.DataSourceString().Computed().Optional().Description("The unique ID for this resource.").Build(),
			},
		},
	}
}

func MapToTeamsDatasourceModel(t *teams.Team) TeamTypeDatasourceModel {
	var team TeamTypeDatasourceModel
	team.CanBeDeleted = types.BoolValue(t.CanBeDeleted)
	team.CanBeRenamed = types.BoolValue(t.CanBeRenamed)
	team.CanChangeMembers = types.BoolValue(t.CanChangeMembers)
	team.CanChangeRoles = types.BoolValue(t.CanChangeRoles)
	team.Description = types.StringValue(t.Description)
	team.ExternalSecurityGroups = MapToExternalSecurityGroupsDatasourceModel(t.ExternalSecurityGroups)
	team.Name = types.StringValue(t.Name)
	team.SpaceId = types.StringValue(t.SpaceID)
	team.Users = basetypes.SetValue(util.FlattenStringList(t.MemberUserIDs))

	team.ID = types.StringValue(t.ID)
	return team
}

func MapToExternalSecurityGroupsDatasourceModel(es []core.NamedReferenceItem) types.List {
	if es == nil || len(es) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: getExternalSecurityGroupsAttrTypes(),
		})
	}

	groups := make([]attr.Value, 0, len(es))
	for _, g := range es {
		group := map[string]attr.Value{
			"display_id_and_name": types.BoolValue(g.DisplayIDAndName),
			"display_name":        types.StringValue(g.DisplayName),
			"id":                  types.StringValue(g.ID),
		}
		groups = append(groups, types.ObjectValueMust(getExternalSecurityGroupsAttrTypes(), group))
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getExternalSecurityGroupsAttrTypes()}, groups)
}

func getExternalSecurityGroupsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_id_and_name": types.BoolType,
		"display_name":        types.StringType,
		"id":                  types.StringType,
	}
}

type TeamTypeDatasourceModel struct {
	CanBeDeleted           types.Bool   `tfsdk:"can_be_deleted"`
	CanBeRenamed           types.Bool   `tfsdk:"can_be_renamed"`
	CanChangeMembers       types.Bool   `tfsdk:"can_change_members"`
	CanChangeRoles         types.Bool   `tfsdk:"can_change_roles"`
	Description            types.String `tfsdk:"description"`
	ExternalSecurityGroups types.List   `tfsdk:"external_security_group"`
	Name                   types.String `tfsdk:"name"`
	SpaceId                types.String `tfsdk:"space_id"`
	Users                  types.Set    `tfsdk:"users"`
	ResourceModel
}

type TeamExternalSecurityGroupTypeDatasourceModel struct {
	DisplayIdAndName types.Bool   `tfsdk:"display_id_and_name"`
	DisplayName      types.String `tfsdk:"display_name"`

	ResourceModel
}

type TeamTypeResourceModel struct {
	UserRole []ScopedUserRoleResourceModel `tfsdk:"user_role"`

	TeamTypeDatasourceModel
}

type ScopedUserRoleResourceModel struct {
	EnvironmentIDs  types.Set    `tfsdk:"environment_ids"`
	ID              types.String `tfsdk:"id"`
	ProjectGroupIDs types.Set    `tfsdk:"project_group_ids"`
	ProjectIDs      types.Set    `tfsdk:"project_ids"`
	SpaceID         types.String `tfsdk:"space_id"`
	TeamID          types.String `tfsdk:"team_id"`
	TenantIDs       types.Set    `tfsdk:"tenant_ids"`
	UserRoleID      types.String `tfsdk:"user_role_id"`
}
