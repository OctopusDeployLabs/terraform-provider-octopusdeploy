package schemas

import (
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const awsElasticContainerRegistryFeedDescription = "aws elastic container registry"

func GetAwsElasticContainerRegistryFeedResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"access_key": resourceSchema.StringAttribute{
			Required:    true,
			Description: "The AWS access key to use when authenticating against Amazon Web Services.",
		},
		"id":                                   GetIdResourceSchema(),
		"name":                                 GetNameResourceSchema(true),
		"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
		"region": resourceSchema.StringAttribute{
			Required:    true,
			Description: "The AWS region where the registry resides.",
		},
		"secret_key": resourceSchema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "The AWS secret key to use when authenticating against Amazon Web Services.",
		},
		"space_id": GetSpaceIdResourceSchema(awsElasticContainerRegistryFeedDescription),
	}
}

type AwsElasticContainerRegistryFeedTypeResourceModel struct {
	AccessKey                         types.String `tfsdk:"access_key"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Region                            types.String `tfsdk:"region"`
	SecretKey                         types.String `tfsdk:"secret_key"`
	SpaceID                           types.String `tfsdk:"space_id"`

	ResourceModel
}
