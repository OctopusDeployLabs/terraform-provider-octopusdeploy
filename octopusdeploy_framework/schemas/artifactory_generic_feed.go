package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const artifactoryGenericFeedDescription = "artifactory generic feed"

type ArtifactoryGenericFeedSchema struct{}

var _ EntitySchema = ArtifactoryGenericFeedSchema{}

func (a ArtifactoryGenericFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a Artifactory Generic feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"feed_uri": resourceSchema.StringAttribute{
				Required: true,
			},
			"id":                                   GetIdResourceSchema(),
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"password":                             GetPasswordResourceSchema(false),
			"space_id":                             GetSpaceIdResourceSchema(artifactoryGenericFeedDescription),
			"username":                             GetUsernameResourceSchema(false),
			"repository": resourceSchema.StringAttribute{
				Computed: false,
				Required: true,
			},
			"layout_regex": resourceSchema.StringAttribute{
				Computed: false,
				Required: false,
				Optional: true,
			},
		},
	}
}

func (a ArtifactoryGenericFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
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
