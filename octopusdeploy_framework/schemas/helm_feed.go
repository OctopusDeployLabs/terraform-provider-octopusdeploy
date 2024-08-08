package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const helmFeedDescription = "helm feed"

func GetHelmFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"feed_uri":                             util.GetFeedUriResourceSchema(),
		"id":                                   util.GetIdResourceSchema(),
		"name":                                 util.GetNameResourceSchema(true),
		"package_acquisition_location_options": util.GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             util.GetPasswordResourceSchema(false),
		"space_id":                             util.GetSpaceIdResourceSchema(helmFeedDescription),
		"username":                             util.GetUsernameResourceSchema(false),
	}
}

type HelmFeedTypeResourceModel struct {
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`

	ResourceModel
}
