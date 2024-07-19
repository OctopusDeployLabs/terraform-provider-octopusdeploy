package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const gitHubRepositoryFeedDescription = "github repository feed"

func GetGitHubRepositoryFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"download_attempts": resourceSchema.Int64Attribute{
			Default:     int64default.StaticInt64(5),
			Description: "The number of times a deployment should attempt to download a package from this feed before failing.",
			Optional:    true,
			Computed:    true,
		},
		"download_retry_backoff_seconds": resourceSchema.Int64Attribute{
			Default:     int64default.StaticInt64(10),
			Description: "The number of seconds to apply as a linear back off between download attempts.",
			Optional:    true,
			Computed:    true,
		},
		"feed_uri": resourceSchema.StringAttribute{
			Required: true,
		},
		"id":   util.GetIdResourceSchema(),
		"name": util.GetNameResourceSchema(true),
		"package_acquisition_location_options": resourceSchema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Optional:    true,
		},
		"password": util.GetPasswordSchema(false, 1000),
		"space_id": util.GetSpaceIdResourceSchema(gitHubRepositoryFeedDescription),
		"username": util.GetUsernameSchema(false, 1000),
	}
}

type GitHubRepositoryFeedTypeResourceModel struct {
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
}
