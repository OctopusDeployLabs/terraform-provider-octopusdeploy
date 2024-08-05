package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const nugetFeedDescription = "nuget feed"

func GetNugetFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"download_attempts":              util.GetDownloadAttemptsResourceSchema(),
		"download_retry_backoff_seconds": util.GetDownloadRetryBackoffSecondsResourceSchema(),
		"feed_uri":                       util.GetFeedUriResourceSchema(),
		"id":                             util.GetIdResourceSchema(),
		"is_enhanced_mode": resourceSchema.BoolAttribute{
			Computed:    true,
			Default:     booldefault.StaticBool(true),
			Description: "This will improve performance of the NuGet feed but may not be supported by some older feeds. Disable if the operation, Create Release does not return the latest version for a package.",
			Optional:    true,
		},
		"name":                                 util.GetNameResourceSchema(true),
		"package_acquisition_location_options": util.GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             util.GetPasswordResourceSchema(false),
		"space_id":                             util.GetSpaceIdResourceSchema(nugetFeedDescription),
		"username":                             util.GetUsernameResourceSchema(false),
	}
}

type NugetFeedTypeResourceModel struct {
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	ID                                types.String `tfsdk:"id"`
	IsEnhancedMode                    types.Bool   `tfsdk:"is_enhanced_mode"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
}
