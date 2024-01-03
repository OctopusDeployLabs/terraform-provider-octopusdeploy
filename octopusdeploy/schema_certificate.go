package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandCertificate(d *schema.ResourceData) *certificates.CertificateResource {
	name := d.Get("name").(string)
	certificateData := core.NewSensitiveValue(d.Get("certificate_data").(string))
	password := core.NewSensitiveValue(d.Get("password").(string))

	certificate := certificates.NewCertificateResource(name, certificateData, password)
	certificate.ID = d.Id()

	if v, ok := d.GetOk("archived"); ok {
		certificate.Archived = v.(string)
	}

	if v, ok := d.GetOk("certificate_data_format"); ok {
		certificate.CertificateDataFormat = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		certificate.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("has_private_key"); ok {
		certificate.HasPrivateKey = v.(bool)
	}

	if v, ok := d.GetOk("is_expired"); ok {
		certificate.IsExpired = v.(bool)
	}

	if v, ok := d.GetOk("issuer_common_name"); ok {
		certificate.IssuerCommonName = v.(string)
	}

	if v, ok := d.GetOk("issuer_distinguished_name"); ok {
		certificate.IssuerDistinguishedName = v.(string)
	}

	if v, ok := d.GetOk("issuer_organization"); ok {
		certificate.IssuerOrganization = v.(string)
	}

	if v, ok := d.GetOk("not_after"); ok {
		certificate.NotAfter = v.(string)
	}

	if v, ok := d.GetOk("not_before"); ok {
		certificate.NotBefore = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		certificate.Notes = v.(string)
	}

	if v, ok := d.GetOk("replaced_by"); ok {
		certificate.ReplacedBy = v.(string)
	}

	if v, ok := d.GetOk("self_signed"); ok {
		certificate.SelfSigned = v.(bool)
	}

	if v, ok := d.GetOk("serial_number"); ok {
		certificate.SerialNumber = v.(string)
	}

	if v, ok := d.GetOk("signature_algorithm_name"); ok {
		certificate.SignatureAlgorithmName = v.(string)
	}

	if v, ok := d.GetOk("subject_alternative_names"); ok {
		certificate.SubjectAlternativeNames = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("subject_common_name"); ok {
		certificate.SubjectCommonName = v.(string)
	}

	if v, ok := d.GetOk("subject_distinguished_name"); ok {
		certificate.SubjectDistinguishedName = v.(string)
	}

	if v, ok := d.GetOk("subject_organization"); ok {
		certificate.SubjectOrganization = v.(string)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		certificate.TenantedDeploymentMode = core.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenants"); ok {
		certificate.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		certificate.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("thumbprint"); ok {
		certificate.Thumbprint = v.(string)
	}

	if v, ok := d.GetOk("version"); ok {
		certificate.Version = v.(int)
	}

	if v, ok := d.GetOk("space_id"); ok {
		certificate.SpaceID = v.(string)
	}

	return certificate
}

func flattenCertificate(certificate *certificates.CertificateResource) map[string]interface{} {
	if certificate == nil {
		return nil
	}

	// NOTE: certificate fields like certificate_data and password are not
	// present here because they are sensitive values are can only be created
	// or updated; never read

	return map[string]interface{}{
		"archived":                          certificate.Archived,
		"certificate_data_format":           certificate.CertificateDataFormat,
		"environments":                      certificate.EnvironmentIDs,
		"has_private_key":                   certificate.HasPrivateKey,
		"id":                                certificate.GetID(),
		"is_expired":                        certificate.IsExpired,
		"issuer_common_name":                certificate.IssuerCommonName,
		"issuer_distinguished_name":         certificate.IssuerDistinguishedName,
		"issuer_organization":               certificate.IssuerOrganization,
		"name":                              certificate.Name,
		"not_after":                         certificate.NotAfter,
		"not_before":                        certificate.NotBefore,
		"notes":                             certificate.Notes,
		"replaced_by":                       certificate.ReplacedBy,
		"self_signed":                       certificate.SelfSigned,
		"serial_number":                     certificate.SerialNumber,
		"signature_algorithm_name":          certificate.SignatureAlgorithmName,
		"subject_alternative_names":         certificate.SubjectAlternativeNames,
		"subject_common_name":               certificate.SubjectCommonName,
		"subject_distinguished_name":        certificate.SubjectDistinguishedName,
		"subject_organization":              certificate.SubjectOrganization,
		"tenanted_deployment_participation": certificate.TenantedDeploymentMode,
		"tenants":                           certificate.TenantIDs,
		"tenant_tags":                       certificate.TenantTags,
		"thumbprint":                        certificate.Thumbprint,
		"version":                           certificate.Version,
	}
}

func getCertificateDataSchema() map[string]*schema.Schema {
	dataSchema := getCertificateSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"archived": getQueryArchived(),
		"certificates": {
			Computed:    true,
			Description: "A list of certificates that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"first_result": getQueryFirstResult(),
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"order_by":     getQueryOrderBy(),
		"partial_name": getQueryPartialName(),
		"search":       getQuerySearch(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"tenant":       getQueryTenant(),
		"space_id":     getSpaceIDSchema(),
	}
}

func getCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"archived": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificate_data": {
			Description:      "The encoded data of the certificate.",
			Required:         true,
			Sensitive:        true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"certificate_data_format": getCertificateDataFormatSchema(),
		"environments":            getEnvironmentsSchema(),
		"has_private_key": {
			Computed:    true,
			Description: "Indicates if the certificate has a private key.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"id": getIDSchema(),
		"is_expired": {
			Computed:    true,
			Description: "Indicates if the certificate has expired.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"issuer_common_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"issuer_distinguished_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"issuer_organization": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": getNameSchema(true),
		"not_after": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"not_before": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"notes": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"password": getPasswordSchema(false),
		"replaced_by": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"self_signed": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"serial_number": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"signature_algorithm_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_alternative_names": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"subject_common_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_distinguished_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_organization": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"thumbprint": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"version": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"space_id": {
			Optional: true,
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}

func setCertificate(ctx context.Context, d *schema.ResourceData, certificate *certificates.CertificateResource) error {
	d.Set("archived", certificate.Archived)
	d.Set("certificate_data_format", certificate.CertificateDataFormat)
	d.Set("has_private_key", certificate.HasPrivateKey)
	d.Set("is_expired", certificate.IsExpired)
	d.Set("issuer_common_name", certificate.IssuerCommonName)
	d.Set("issuer_distinguished_name", certificate.IssuerDistinguishedName)
	d.Set("issuer_organization", certificate.IssuerOrganization)
	d.Set("name", certificate.Name)
	d.Set("not_after", certificate.NotAfter)
	d.Set("not_before", certificate.NotBefore)
	d.Set("notes", certificate.Notes)
	d.Set("replaced_by", certificate.ReplacedBy)
	d.Set("serial_number", certificate.SerialNumber)
	d.Set("signature_algorithm_name", certificate.SignatureAlgorithmName)
	d.Set("subject_common_name", certificate.SubjectCommonName)
	d.Set("subject_distinguished_name", certificate.SubjectDistinguishedName)
	d.Set("subject_organization", certificate.SubjectOrganization)
	d.Set("self_signed", certificate.SelfSigned)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentMode)
	d.Set("thumbprint", certificate.Thumbprint)
	d.Set("version", certificate.Version)

	if err := d.Set("environments", certificate.EnvironmentIDs); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("subject_alternative_names", certificate.SubjectAlternativeNames); err != nil {
		return fmt.Errorf("error setting subject_alternative_names: %s", err)
	}

	if err := d.Set("tenants", certificate.TenantIDs); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", certificate.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	return nil
}
