package schemas

import (
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const helmFeedDescription = "helm feed"

func GetHelmFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"feed_uri":                             GetFeedUriResourceSchema(),
		"id":                                   GetIdResourceSchema(),
		"name":                                 GetNameResourceSchema(true),
		"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             GetPasswordResourceSchema(false),
		"space_id":                             GetSpaceIdResourceSchema(helmFeedDescription),
		"username":                             GetUsernameResourceSchema(false),
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
