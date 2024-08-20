package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const gitHubRepositoryFeedDescription = "github repository feed"

type GitHubRepositoryFeedSchema struct{}

var _ EntitySchema = GitHubRepositoryFeedSchema{}

func (g GitHubRepositoryFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a GitHub repository feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"download_attempts":                    GetDownloadAttemptsResourceSchema(),
			"download_retry_backoff_seconds":       GetDownloadRetryBackoffSecondsResourceSchema(),
			"feed_uri":                             GetFeedUriResourceSchema(),
			"id":                                   GetIdResourceSchema(),
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"password":                             GetPasswordResourceSchema(false),
			"space_id":                             GetSpaceIdResourceSchema(gitHubRepositoryFeedDescription),
			"username":                             GetUsernameResourceSchema(false),
		},
	}
}

func (g GitHubRepositoryFeedSchema) GetDatasourceSchemaAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{}
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
