package schemas

import (
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const dockerContainerRegistryFeedDescription = "docker container registry feed"

func GetDockerContainerRegistryFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"api_version": resourceSchema.StringAttribute{
			Optional: true,
		},
		"feed_uri":                             GetFeedUriResourceSchema(),
		"id":                                   GetIdResourceSchema(),
		"name":                                 GetNameResourceSchema(true),
		"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             GetPasswordResourceSchema(false),
		"space_id":                             GetSpaceIdResourceSchema(dockerContainerRegistryFeedDescription),
		"username":                             GetUsernameResourceSchema(false),
		"registry_path": resourceSchema.StringAttribute{
			Optional: true,
		},
	}
}

type DockerContainerRegistryFeedTypeResourceModel struct {
	APIVersion                        types.String `tfsdk:"api_version"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
	RegistryPath                      types.String `tfsdk:"registry_path"`

	ResourceModel
}
