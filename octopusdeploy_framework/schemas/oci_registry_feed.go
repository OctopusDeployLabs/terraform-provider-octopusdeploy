package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ociRegistryFeedDescription = "OCI registry"

type OCIRegistryFeedSchema struct{}

func (m OCIRegistryFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a OCI Registry feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"feed_uri": GetFeedUriResourceSchema(),
			"id":       GetIdResourceSchema(),
			"name":     GetNameResourceSchema(true),
			"password": GetPasswordResourceSchema(false),
			"space_id": GetSpaceIdResourceSchema(ociRegistryFeedDescription),
			"username": GetUsernameResourceSchema(false),
		},
	}
}

func (m OCIRegistryFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = OCIRegistryFeedSchema{}

type OCIRegistryFeedTypeResourceModel struct {
	FeedUri  types.String `tfsdk:"feed_uri"`
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	SpaceID  types.String `tfsdk:"space_id"`
	Username types.String `tfsdk:"username"`

	ResourceModel
}
