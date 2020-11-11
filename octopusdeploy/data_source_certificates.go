package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificatesRead,
		Schema:      getCertificateDataSchema(),
	}
}

func dataSourceCertificatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.CertificatesQuery{
		Archived:    d.Get("archived").(string),
		FirstResult: d.Get("first_result").(string),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		OrderBy:     d.Get("order_by").(string),
		PartialName: d.Get("partial_name").(string),
		Search:      d.Get("search").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
		Tenant:      d.Get("tenant").(string),
	}

	client := m.(*octopusdeploy.Client)
	certificates, err := client.Certificates.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCertificates := []interface{}{}
	for _, certificate := range certificates.Items {
		flattenedCertificate := map[string]interface{}{
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
			"serial_number":                     certificate.SerialNumber,
			"signature_algorithm_name":          certificate.SignatureAlgorithmName,
			"subject_alternative_names":         certificate.SubjectAlternativeNames,
			"subject_common_name":               certificate.SubjectCommonName,
			"subject_distinguished_name":        certificate.SubjectDistinguishedName,
			"subject_organization":              certificate.SubjectOrganization,
			"self_signed":                       certificate.SelfSigned,
			"tenanted_deployment_participation": certificate.TenantedDeploymentMode,
			"tenants":                           certificate.TenantIDs,
			"tenant_tags":                       certificate.TenantTags,
			"thumbprint":                        certificate.Thumbprint,
			"version":                           certificate.Version,
		}
		flattenedCertificates = append(flattenedCertificates, flattenedCertificate)
	}

	d.Set("certificates", flattenedCertificates)
	d.SetId("Certificates " + time.Now().UTC().String())

	return nil
}
