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
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
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
