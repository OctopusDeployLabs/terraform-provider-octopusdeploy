package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
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
		"id":   util.GetIdResourceSchema(),
		"name": util.GetNameResourceSchema(true),
		"package_acquisition_location_options": resourceSchema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Optional:    true,
		},
		"region": resourceSchema.StringAttribute{
			Required:    true,
			Description: "The AWS region where the registry resides.",
		},
		"secret_key": resourceSchema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "The AWS secret key to use when authenticating against Amazon Web Services.",
		},
		"space_id": util.GetSpaceIdResourceSchema(awsElasticContainerRegistryFeedDescription),
	}
}

type AwsElasticContainerRegistryFeedTypeResourceModel struct {
	AccessKey                         types.String `tfsdk:"access_key"`
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List   `tfsdk:"package_acquisition_location_options"`
	Region                            types.String `tfsdk:"region"`
	SecretKey                         types.String `tfsdk:"secret_key"`
	SpaceID                           types.String `tfsdk:"space_id"`
}
