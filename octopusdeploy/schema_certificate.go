package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandCertificate(d *schema.ResourceData) *octopusdeploy.CertificateResource {
	name := d.Get("name").(string)
	certificateData := octopusdeploy.NewSensitiveValue(d.Get("certificate_data").(string))
	password := octopusdeploy.NewSensitiveValue(d.Get("password").(string))

	certificate := octopusdeploy.NewCertificateResource(name, certificateData, password)
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
		certificate.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
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

	return certificate
}

func flattenCertificate(certificate *octopusdeploy.CertificateResource) map[string]interface{} {
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

func setCertificate(ctx context.Context, d *schema.ResourceData, certificate *octopusdeploy.CertificateResource) {
	// NOTE: certificate fields like certificate_data and password are not
	// present here because they are sensitive values are can only be created
	// or updated; never read

	d.Set("archived", certificate.Archived)
	d.Set("certificate_data_format", certificate.CertificateDataFormat)
	d.Set("environments", certificate.EnvironmentIDs)
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
	d.Set("subject_alternative_names", certificate.SubjectAlternativeNames)
	d.Set("subject_common_name", certificate.SubjectCommonName)
	d.Set("subject_distinguished_name", certificate.SubjectDistinguishedName)
	d.Set("subject_organization", certificate.SubjectOrganization)
	d.Set("self_signed", certificate.SelfSigned)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentMode)
	d.Set("tenants", certificate.TenantIDs)
	d.Set("tenant_tags", certificate.TenantTags)
	d.Set("thumbprint", certificate.Thumbprint)
	d.Set("version", certificate.Version)

	d.SetId(certificate.GetID())
}

func getCertificateDataSchema() map[string]*schema.Schema {
	certificateSchema := getCertificateSchema()
	for _, field := range certificateSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

	return map[string]*schema.Schema{
		"archived": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificates": {
			Computed: true,
			Elem:     &schema.Resource{Schema: certificateSchema},
			Type:     schema.TypeList,
		},
		"first_result": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"order_by": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"search": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"take": {
			Default:  1,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"tenant": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}

func getCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"archived": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificate_data": {
			Required:     true,
			Sensitive:    true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"certificate_data_format": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"has_private_key": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_expired": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"issuer_common_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"issuer_distinguished_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"issuer_organization": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"not_after": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"not_before": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"notes": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"password": {
			Required:     true,
			Sensitive:    true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"replaced_by": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"self_signed": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"serial_number": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"signature_algorithm_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_alternative_names": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subject_common_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_distinguished_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"subject_organization": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tenant_tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"version": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
