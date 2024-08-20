package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const nugetFeedDescription = "nuget feed"

type NugetFeedSchema struct{}

func (n NugetFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"download_attempts":              GetDownloadAttemptsResourceSchema(),
			"download_retry_backoff_seconds": GetDownloadRetryBackoffSecondsResourceSchema(),
			"feed_uri":                       GetFeedUriResourceSchema(),
			"id":                             GetIdResourceSchema(),
			"is_enhanced_mode": resourceSchema.BoolAttribute{
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "This will improve performance of the NuGet feed but may not be supported by some older feeds. Disable if the operation, Create Release does not return the latest version for a package.",
				Optional:    true,
			},
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"password":                             GetPasswordResourceSchema(false),
			"space_id":                             GetSpaceIdResourceSchema(nugetFeedDescription),
			"username":                             GetUsernameResourceSchema(false),
		},
		Description: "This resource manages a Nuget feed in Octopus Deploy.",
	}
}

func (n NugetFeedSchema) GetDatasource() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = NugetFeedSchema{}

type NugetFeedTypeResourceModel struct {
	DownloadAttempts                  types.Int64  `tfsdk:"download_attempts"`
	DownloadRetryBackoffSeconds       types.Int64  `tfsdk:"download_retry_backoff_seconds"`
	FeedUri                           types.String `tfsdk:"feed_uri"`
	IsEnhancedMode                    types.Bool   `tfsdk:"is_enhanced_mode"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`

	ResourceModel
}
