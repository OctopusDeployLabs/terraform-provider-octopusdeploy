package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const awsElasticContainerRegistryFeedDescription = "aws elastic container registry"

type AwsElasticContainerRegistrySchema struct{}

var _ EntitySchema = AwsElasticContainerRegistrySchema{}

func (a AwsElasticContainerRegistrySchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (a AwsElasticContainerRegistrySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages an AWS Elastic Container Registry in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"access_key": resourceSchema.StringAttribute{
				Optional:    true,
				Description: "The AWS access key to use when authenticating against Amazon Web Services.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id":                                   GetIdResourceSchema(),
			"name":                                 GetNameResourceSchema(true),
			"package_acquisition_location_options": GetPackageAcquisitionLocationOptionsResourceSchema(),
			"region": resourceSchema.StringAttribute{
				Required:    true,
				Description: "The AWS region where the registry resides.",
			},
			"secret_key": resourceSchema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The AWS secret key to use when authenticating against Amazon Web Services.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": GetSpaceIdResourceSchema(awsElasticContainerRegistryFeedDescription),
			"oidc_authentication": resourceSchema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceSchema.Attribute{
					"session_duration": resourceSchema.StringAttribute{
						Description: "Assumed role session duration (in seconds)",
						Optional:    true,
						Computed:    true,
					},
					"audience": resourceSchema.StringAttribute{
						Description: "Audience to use when authenticating against Amazon Web Services.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"role_arn": resourceSchema.StringAttribute{
						Description: "The Amazon Resource Name (ARN) of the role that the caller is assuming.",
						Computed:    true,
						Optional:    true,
						Default:     stringdefault.StaticString(""),
					},
					"subject_keys": GetOidcSubjectKeysSchema("Keys to include in a deployment or runbook. Valid options are `space`, `feed`.", false),
				},
			},
		},
	}
}

type AwsElasticContainerRegistryFeedTypeResourceModel struct {
	AccessKey                         types.String                        `tfsdk:"access_key"`
	Name                              types.String                        `tfsdk:"name"`
	PackageAcquisitionLocationOptions types.List                          `tfsdk:"package_acquisition_location_options"`
	Region                            types.String                        `tfsdk:"region"`
	SecretKey                         types.String                        `tfsdk:"secret_key"`
	SpaceID                           types.String                        `tfsdk:"space_id"`
	OidcAuthentication                *EcrOidcAuthenticationResourceModel `tfsdk:"oidc_authentication"`

	ResourceModel
}

type EcrOidcAuthenticationResourceModel struct {
	SessionDuration types.String `tfsdk:"session_duration"`
	Audience        types.String `tfsdk:"audience"`
	RoleArn         types.String `tfsdk:"role_arn"`
	SubjectKey      types.List   `tfsdk:"subject_keys"`
}
