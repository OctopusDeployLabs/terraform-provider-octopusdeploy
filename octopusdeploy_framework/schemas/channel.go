package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ChannelResourceDescription = "channel"

type ChannelSchema struct{}

type ChannelModel struct {
	Description types.String `tfsdk:"description"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	LifecycleId types.String `tfsdk:"lifecycle_id"`
	Name        types.String `tfsdk:"name"`
	ProjectId   types.String `tfsdk:"project_id"`
	Rule        types.List   `tfsdk:"rule"`
	SpaceId     types.String `tfsdk:"space_id"`
	TenantTags  types.List   `tfsdk:"tenant_tags"`

	ResourceModel
}

func (c ChannelSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(ChannelResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"description": GetDescriptionResourceSchema(ChannelResourceDescription),
			"is_default": resourceSchema.BoolAttribute{
				Description: "Indicates whether this is the default channel for the associated project.",
				Optional:    true,
			},
			"lifecycle_id": resourceSchema.StringAttribute{
				Description: "The lifecycle ID associated with this channel.",
				Optional:    true,
			},
			"name": GetNameResourceSchema(true),
			"project_id": resourceSchema.StringAttribute{
				Description: "The project ID associated with this channel.",
				Required:    true,
			},
			"rule": resourceSchema.ListNestedAttribute{
				Description: "A list of rules associated with this channel.",
				Optional:    true,
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"action_package": resourceSchema.ListNestedAttribute{
							Required: true,
							NestedObject: resourceSchema.NestedAttributeObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": resourceSchema.StringAttribute{
										Optional: true,
									},
									"package_reference": resourceSchema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
						"id": resourceSchema.StringAttribute{
							Description: "The ID associated with this channel rule.",
							Computed:    true,
							Optional:    true,
						},
						"tag": resourceSchema.StringAttribute{
							Optional: true,
						},
						"version_range": resourceSchema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"space_id": GetSpaceIdResourceSchema(ChannelResourceDescription),
			"tenant_tags": resourceSchema.ListAttribute{
				Description: "A list of tenant tags associated with this channel.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}
