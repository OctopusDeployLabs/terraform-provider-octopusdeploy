package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const gitHubRepositoryFeedDescription = "github repository feed"

func GetGitHubRepositoryFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"download_attempts":                    util.GetDownloadAttemptsResourceSchema(),
		"download_retry_backoff_seconds":       util.GetDownloadRetryBackoffSecondsResourceSchema(),
		"feed_uri":                             util.GetFeedUriResourceSchema(),
		"id":                                   util.GetIdResourceSchema(),
		"name":                                 util.GetNameResourceSchema(true),
		"package_acquisition_location_options": util.GetPackageAcquisitionLocationOptionsResourceSchema(),
		"password":                             util.GetPasswordResourceSchema(false),
		"space_id":                             util.GetSpaceIdResourceSchema(gitHubRepositoryFeedDescription),
		"username":                             util.GetUsernameResourceSchema(false),
	}
}

type GitHubRepositoryFeedTypeResourceModel struct {
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`

	ResourceModel
}
