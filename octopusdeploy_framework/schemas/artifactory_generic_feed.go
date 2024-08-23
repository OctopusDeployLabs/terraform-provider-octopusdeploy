package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ArtifactoryGenericFeedSchema struct{}

func GetArtifactoryGenericFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"feed_uri": resourceSchema.StringAttribute{
			Required: true,
		},
		"id":                                   util.GetIdResourceSchema(),
		"name":                                 util.GetNameResourceSchema(true),
		"package_acquisition_location_options": util.GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             util.GetPasswordResourceSchema(false),
		"space_id":                             util.GetSpaceIdResourceSchema(artifactoryGenericFeedDescription),
		"username":                             util.GetUsernameResourceSchema(false),
		"repository": resourceSchema.StringAttribute{
			Computed: false,
			Required: true,
		},
		"layout_regex": resourceSchema.StringAttribute{
			Computed: false,
			Required: false,
			Optional: true,
		},
	}
}

type ArtifactoryGenericFeedTypeResourceModel struct {
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
	Repository                        types.String `tfsdk:"repository"`
	LayoutRegex                       types.String `tfsdk:"layout_regex"`

	ResourceModel
}
