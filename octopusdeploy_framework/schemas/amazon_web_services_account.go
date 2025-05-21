package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const AmazonWebServicesAccountResourceDescription = "AWS account"

type AmazonWebServicesAccountSchema struct{}

type AmazonWebServicesAccountModel struct {
	AccessKey                       types.String `tfsdk:"access_key"`
	Description                     types.String `tfsdk:"description"`
	Environments                    types.List   `tfsdk:"environments"`
	Name                            types.String `tfsdk:"name"`
	SecretKey                       types.String `tfsdk:"secret_key"`
	SpaceId                         types.String `tfsdk:"space_id"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	Tenants                         types.List   `tfsdk:"tenants"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`

	ResourceModel
}

func (a AmazonWebServicesAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(AmazonWebServicesAccountResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"access_key": resourceSchema.StringAttribute{
				Description: "The access key associated with this AWS account.",
				Required:    true,
			},
			"description": GetDescriptionResourceSchema(AmazonWebServicesAccountResourceDescription),
			"environments": resourceSchema.ListAttribute{
				Description: "A list of environment IDs associated with this AWS account.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"name": GetNameResourceSchema(true),
			"secret_key": resourceSchema.StringAttribute{
				Description: "The secret key associated with this AWS account.",
				Sensitive:   true,
				Required:    true,
			},
			"space_id": GetSpaceIdResourceSchema(AmazonWebServicesAccountResourceDescription),
			"tenanted_deployment_participation": resourceSchema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Untenanted", "TenantedOrUntenanted", "Tenanted"),
				},
			},
			"tenants": resourceSchema.ListAttribute{
				Description: "A list of tenant IDs associated with this AWS account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tenant_tags": resourceSchema.ListAttribute{
				Description: "A list of tenant tags associated with this AWS account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}
