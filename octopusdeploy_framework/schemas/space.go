package schemas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const spaceDescription = "space"

type SpaceSchema struct{}

var _ EntitySchema = SpaceSchema{}

type SpaceModel struct {
	Name                     types.String `tfsdk:"name"`
	Slug                     types.String `tfsdk:"slug"`
	Description              types.String `tfsdk:"description"`
	IsDefault                types.Bool   `tfsdk:"is_default"`
	SpaceManagersTeams       types.Set    `tfsdk:"space_managers_teams"`
	SpaceManagersTeamMembers types.Set    `tfsdk:"space_managers_team_members"`
	IsTaskQueueStopped       types.Bool   `tfsdk:"is_task_queue_stopped"`

	ResourceModel
}

func (s SpaceSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages spaces in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"description": GetDescriptionResourceSchema(spaceDescription),
			"name":        GetNameResourceSchema(true),
			"slug":        GetSlugResourceSchema(spaceDescription),
			"space_managers_teams": resourceSchema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of team IDs designated to be managers of this space.",
				Optional:    true,
				Computed:    true,
			},
			"space_managers_team_members": resourceSchema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of user IDs designated to be managers of this space.",
				Optional:    true,
				Computed:    true,
			},
			"is_task_queue_stopped": resourceSchema.BoolAttribute{
				Description: "Specifies the status of the task queue for this space.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_default": resourceSchema.BoolAttribute{
				Description: "Specifies if this space is the default space in Octopus.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (s SpaceSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about an existing space.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":          GetIdDatasourceSchema(true),
			"description": GetReadonlyDescriptionDatasourceSchema(spaceDescription),
			"name": datasourceSchema.StringAttribute{
				Description: fmt.Sprintf("The name of this resource, no more than %d characters long", 20),
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 20),
				},
				Computed: true,
				Optional: true,
			},
			"slug": GetSlugDatasourceSchema(spaceDescription, true),
			"space_managers_teams": datasourceSchema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of team IDs designated to be managers of this space.",
				Computed:    true,
			},
			"space_managers_team_members": datasourceSchema.SetAttribute{
				ElementType: types.StringType,
				Description: "A list of user IDs designated to be managers of this space.",
				Computed:    true,
			},
			"is_task_queue_stopped": datasourceSchema.BoolAttribute{
				Description: "Specifies the status of the task queue for this space.",
				Computed:    true,
			},
			"is_default": datasourceSchema.BoolAttribute{
				Description: "Specifies if this space is the default space in Octopus.",
				Computed:    true,
			},
		},
	}
}

func GetSpaceTypeAttributes() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":                          types.StringType,
		"name":                        types.StringType,
		"slug":                        types.StringType,
		"description":                 types.StringType,
		"is_default":                  types.BoolType,
		"space_managers_teams":        types.SetType{ElemType: types.StringType},
		"space_managers_team_members": types.SetType{ElemType: types.StringType},
		"is_task_queue_stopped":       types.BoolType}}
}
