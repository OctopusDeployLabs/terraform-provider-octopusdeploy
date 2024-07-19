package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const artifactoryGenericFeedDescription = "artifactory generic feed"

func GetArtifactoryGenericFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
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
		"password": util.GetPasswordResourceSchema(false),
		"space_id": util.GetSpaceIdResourceSchema(helmFeedDescription),
		"username": util.GetUsernameResourceSchema(false),
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
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Password                          types.String `tfsdk:"password"`
	SpaceID                           types.String `tfsdk:"space_id"`
	Username                          types.String `tfsdk:"username"`
	Repository                        types.String `tfsdk:"repository"`
	LayoutRegex                       types.String `tfsdk:"layout_regex"`
}
