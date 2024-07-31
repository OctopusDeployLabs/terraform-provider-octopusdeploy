package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const dockerContainerRegistryFeedDescription = "docker container registry feed"

func GetDockerContainerRegistryFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"api_version": resourceSchema.StringAttribute{
			Optional: true,
		},
		"feed_uri":                             util.GetFeedUriResourceSchema(),
		"id":                                   util.GetIdResourceSchema(),
		"name":                                 util.GetNameResourceSchema(true),
		"package_acquisition_location_options": util.GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             util.GetPasswordResourceSchema(false),
		"space_id":                             util.GetSpaceIdResourceSchema(mavenFeedDescription),
		"username":                             util.GetUsernameResourceSchema(false),
	}
}

type DockerContainerRegistryFeedTypeResourceModel struct {
	ApiVersion                        types.String `tfsdk:"api_version"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
}
