package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GoogleContainerRegistryFeedSchema struct{}

var _ EntitySchema = GoogleContainerRegistryFeedSchema{}

func (d GoogleContainerRegistryFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a Google Container Registry feed in Octopus Deploy (alias for Docker Container Registry feed)",
		Attributes: map[string]resourceSchema.Attribute{
			"api_version": resourceSchema.StringAttribute{
				Optional: true,
			},
			"feed_uri": GetFeedUriResourceSchema(),
			"id":       GetIdResourceSchema(),
			"name":     GetNameResourceSchema(true),
			"password": GetPasswordResourceSchema(false),
			"space_id": GetSpaceIdResourceSchema("Google container registry feed"),
			"username": GetUsernameResourceSchema(false),
			"registry_path": resourceSchema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (d GoogleContainerRegistryFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type GoogleContainerRegistryFeedTypeResourceModel struct {
	APIVersion   types.String `tfsdk:"api_version"`
	FeedUri      types.String `tfsdk:"feed_uri"`
	Name         types.String `tfsdk:"name"`
	Password     types.String `tfsdk:"password"`
	SpaceID      types.String `tfsdk:"space_id"`
	Username     types.String `tfsdk:"username"`
	RegistryPath types.String `tfsdk:"registry_path"`

	ResourceModel
}
