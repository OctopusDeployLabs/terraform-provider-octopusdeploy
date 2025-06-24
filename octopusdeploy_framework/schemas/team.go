package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ EntitySchema = TeamSchema{}

type TeamSchema struct{}

func (l TeamSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages teams in Octopus Deploy.",
		Attributes:  map[string]resourceSchema.Attribute{},
		Blocks:      map[string]resourceSchema.Block{},
	}
}

func (l TeamSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing teams.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":             util.DataSourceString().Computed().Description("An auto-generated identifier that includes the timestamp when this data source was last modified.").Build(),
			"ids":            util.DataSourceList(types.StringType).Optional().Description("A filter to search by a list of IDs.").Build(),
			"include_system": util.DataSourceBool().Optional().Description("A filter to include system teams.").Build(),
			"partial_name":   util.DataSourceString().Optional().Description("A filter to search by the partial match of a name.").Build(),
			"spaces":         util.DataSourceList(types.StringType).Optional().Description("A filter to search by a list of space IDs.").Build(),
			"skip":           util.DataSourceInt64().Optional().Description("A filter to specify the number of items to skip in the response.").Build(),
			"take":           util.DataSourceInt64().Optional().Description("A filter to specify the number of items to take (or return) in the response.").Build(),
			"teams":          getDatasourceTeamsAttributes(),
		},
	}
}

func getDatasourceTeamsAttributes() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed:    true,
		Description: "A list of teams that match the filter(s).",
		Optional:    false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"can_be_deleted":          util.DataSourceBool().Computed().Build(),
				"can_be_renamed":          util.DataSourceBool().Computed().Build(),
				"can_change_members":      util.DataSourceBool().Computed().Build(),
				"can_change_roles":        util.DataSourceBool().Computed().Build(),
				"description":             util.DataSourceString().Computed().Description("The user-friendly description of this team.").Build(),
				"external_security_group": getDatasourceExternalSecurityGroupsAttributes(),
				"id":                      util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
				"name":                    util.DataSourceString().Computed().Description("The name of this team.").Build(),
				"space_id":                util.DataSourceString().Computed().Description("The space associated with this team.").Build(),
				"users":                   util.DataSourceSet(types.StringType).Computed().Description("A list of user IDs designated to be members of this team.").Build(),
			},
		},
	}
}

func getDatasourceExternalSecurityGroupsAttributes() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"display_id_and_name": util.DataSourceBool().Computed().Build(),
				"display_name":        util.DataSourceString().Computed().Build(),
				"id":                  util.DataSourceString().Computed().Description("The unique ID for this resource.").Build(),
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
