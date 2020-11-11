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

	if v, ok := d.GetOk("certificate_data"); ok {
		certificate.CertificateData = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("has_private_key"); ok {
		certificate.HasPrivateKey = v.(bool)
	}

	if v, ok := d.GetOk("environments"); ok {
		certificate.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("name"); ok {
		certificate.Name = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		certificate.Notes = v.(string)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		certificate.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		certificate.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return certificate
}

func flattenCertificate(ctx context.Context, d *schema.ResourceData, certificate *octopusdeploy.CertificateResource) {
	d.Set("archived", certificate.Archived)

	if certificate.CertificateData != nil {
		d.Set("certificate_data", certificate.CertificateData.NewValue)
	}

	d.Set("has_private_key", certificate.HasPrivateKey)
	d.Set("environments", certificate.EnvironmentIDs)
	d.Set("name", certificate.Name)
	d.Set("notes", certificate.Notes)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentMode)
	d.Set("tenants", certificate.TenantIDs)
	d.Set("tenant_tags", certificate.TenantTags)

	d.SetId(certificate.GetID())
}

func getCertificateDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"archived": {
			Optional: true,
			Type:     schema.TypeString,
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
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"tenant": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"certificates": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"archived": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"certificate_data_format": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"environments": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"has_private_key": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"is_expired": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"issuer_common_name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"issuer_distinguished_name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"issuer_organization": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"not_after": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"not_before": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"notes": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"replaced_by": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"serial_number": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"signature_algorithm_name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"subject_alternative_names": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"subject_common_name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"subject_distinguished_name": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"subject_organization": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"self_signed": &schema.Schema{
						Optional: true,
						Type:     schema.TypeBool,
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
					"thumbprint": &schema.Schema{
						Optional: true,
						Type:     schema.TypeString,
					},
					"version": &schema.Schema{
						Optional: true,
						Type:     schema.TypeInt,
					},
				},
			},
			Type: schema.TypeList,
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
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"has_private_key": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"notes": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tenant_tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
