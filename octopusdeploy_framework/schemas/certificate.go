package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CertificateSchema struct{}

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
