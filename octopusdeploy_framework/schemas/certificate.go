package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CertificateSchema struct{}

type CertificateModel struct {
	Name                     types.String `tfsdk:"name"`
	Archived                 types.String `tfsdk:"archived"`
	CertificateData          types.String `tfsdk:"certificate_data"`
	CertificateDataFormat    types.String `tfsdk:"certificate_data_format"`
	EnvironmentIDs           types.List   `tfsdk:"environments"`
	HasPrivateKey            types.Bool   `tfsdk:"has_private_key"`
	IsExpired                types.Bool   `tfsdk:"is_expired"`
	IssuerCommonName         types.String `tfsdk:"issuer_common_name"`
	IssuerDistinguishedName  types.String `tfsdk:"issuer_distinguished_name"`
	IssuerOrganization       types.String `tfsdk:"issuer_organization"`
	NotAfter                 types.String `tfsdk:"not_after"`
	NotBefore                types.String `tfsdk:"not_before"`
	Notes                    types.String `tfsdk:"notes"`
	Password                 types.String `tfsdk:"password"`
	ReplacedBy               types.String `tfsdk:"replaced_by"`
	SelfSigned               types.Bool   `tfsdk:"self_signed"`
	SerialNumber             types.String `tfsdk:"serial_number"`
	SignatureAlgorithmName   types.String `tfsdk:"signature_algorithm_name"`
	SpaceID                  types.String `tfsdk:"space_id"`
	SubjectAlternativeNames  types.List   `tfsdk:"subject_alternative_names"`
	SubjectCommonName        types.String `tfsdk:"subject_common_name"`
	SubjectDistinguishedName types.String `tfsdk:"subject_distinguished_name"`
	SubjectOrganization      types.String `tfsdk:"subject_organization"`
	TenantedDeploymentMode   types.String `tfsdk:"tenanted_deployment_participation"`
	TenantIDs                types.List   `tfsdk:"tenants"`
	TenantTags               types.List   `tfsdk:"tenant_tags"`
	Thumbprint               types.String `tfsdk:"thumbprint"`
	Version                  types.Int64  `tfsdk:"version"`

	ResourceModel
}

func (c CertificateSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages certificates in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"archived": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"certificate_data": resourceSchema.StringAttribute{
				Description: "The encoded data of the certificate.",
				Required:    true,
				Sensitive:   true,
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"certificate_data_format": getCertificateDataFormatResourceSchema(),
			"environments":            getEnvironmentsResourceSchema(),
			"has_private_key": resourceSchema.BoolAttribute{
				Description: "Indicates if the certificate has a private key.",
				Computed:    true,
				Optional:    true,
			},
			"id": GetIdResourceSchema(),
			"is_expired": resourceSchema.BoolAttribute{
				Description: "Indicates if the certificate has expired.",
				Computed:    true,
				Optional:    true,
			},
			"issuer_common_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"issuer_distinguished_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"issuer_organization": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name": GetNameResourceSchema(true),
			"not_after": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"not_before": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"notes": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"password": GetPasswordResourceSchema(false),
			"replaced_by": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"self_signed": resourceSchema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"serial_number": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"signature_algorithm_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"subject_alternative_names": resourceSchema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"subject_common_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"subject_distinguished_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"subject_organization": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
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
				Description: "A list of tenant IDs associated with this certificate.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tenant_tags": resourceSchema.ListAttribute{
				Description: "A list of tenant tags associated with this certificate.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"thumbprint": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"version": resourceSchema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"space_id": resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}
