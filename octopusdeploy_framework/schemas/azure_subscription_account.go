package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const AzureSubscriptionAccountDescription = "Azure subscription account"

type AzureSubscriptionAccountSchema struct{}

type AzureSubscriptionAccountModel struct {
	AzureEnvironment                types.String `tfsdk:"azure_environment"`
	Certificate                     types.String `tfsdk:"certificate"`
	CertificateThumbprint           types.String `tfsdk:"certificate_thumbprint"`
	Description                     types.String `tfsdk:"description"`
	Environments                    types.List   `tfsdk:"environments"`
	ManagementEndpoint              types.String `tfsdk:"management_endpoint"`
	Name                            types.String `tfsdk:"name"`
	SpaceID                         types.String `tfsdk:"space_id"`
	StorageEndpointSuffix           types.String `tfsdk:"storage_endpoint_suffix"`
	SubscriptionID                  types.String `tfsdk:"subscription_id"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	Tenants                         types.List   `tfsdk:"tenants"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`

	ResourceModel
}

func (a AzureSubscriptionAccountSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(AzureSubscriptionAccountDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"azure_environment": resourceSchema.StringAttribute{
				Description: "The Azure environment associated with this Azure subscription account. Valid Azure environments are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, or `AzureUSGovernment`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AzureCloud",
						"AzureChinaCloud",
						"AzureGermanCloud",
						"AzureUSGovernment",
					),
				},
			},
			"certificate": resourceSchema.StringAttribute{
				Description: "TODO",
				Optional:    true,
				Sensitive:   true,
			},
			"certificate_thumbprint": resourceSchema.StringAttribute{
				Description: "TODO",
				Optional:    true,
				Sensitive:   true,
			},
			"description": GetDescriptionResourceSchema(AzureSubscriptionAccountDescription),
			"environments": resourceSchema.ListAttribute{
				Description: "A list of environment IDs associated with this Azure subscription account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"management_endpoint": resourceSchema.StringAttribute{
				Description: "TODO",
				// TODO: add equivalent of below RequiredWith
				// RequiredWith: []string{"azure_environment"}
			},
			"name":     GetNameResourceSchema(true),
			"space_id": GetSpaceIdResourceSchema(AzureSubscriptionAccountDescription),
			"storage_endpoint_suffix": resourceSchema.StringAttribute{
				Description: "The storage endpoint suffix associated with this Azure subscription account.",
				Optional:    true,
				// TODO: add equivalent of below RequiredWith
				// RequiredWith: []string{"azure_environment"},
			},
			"subscription_id": resourceSchema.StringAttribute{
				Description: "The subscription ID of this resource.",
				Required:    true,
				// TODO: add UUID validator
			},
			"tenanted_deployment_participation": resourceSchema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Untenanted",
						"TenantedOrUntenanted",
						"Tenanted",
					),
				},
			},
			"tenants": resourceSchema.ListAttribute{
				Description: "A list of tenant IDs associated with this Azure subscription account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tenant_tags": resourceSchema.ListAttribute{
				Description: "A list of tenant tags associated with this Azure subscription account.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}
