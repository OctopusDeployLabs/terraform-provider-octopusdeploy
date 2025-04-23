package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const ChannelResourceDescription = "channel"

type ChannelSchema struct{}

type ChannelModel struct {
	IsDefault   bool   `tfsdk:"is_default"`
	LifecycleId string `tfsdk:"lifecycle_id"`
	Name        string `tfsdk:"name"`
	ProjectId   string `tfsdk:"project_id"`
	SpaceId     string `tfsdk:"space_id"`

	ResourceModel
}

func (c ChannelSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(ChannelResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id":         GetIdResourceSchema(),
			"is_default": resourceSchema.BoolAttribute{},
			"lifecycle_id": resourceSchema.StringAttribute{
				Description: "The lifecycle ID associated with the channel.",
				Optional:    true,
			},
			"name": GetNameResourceSchema(true),
			"project_id": resourceSchema.StringAttribute{
				Description: "The project ID associated with the channel.",
				Required:    true,
			},
			"space_id": GetSpaceIdResourceSchema(ChannelResourceDescription),
		},
	}
}
