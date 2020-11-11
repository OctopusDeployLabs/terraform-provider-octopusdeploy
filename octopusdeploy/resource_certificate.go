package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		DeleteContext: resourceCertificateDelete,
		Importer:      getImporter(),
		ReadContext:   resourceCertificateRead,
		Schema:        getCertificateSchema(),
		UpdateContext: resourceCertificateUpdate,
	}
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate := expandCertificate(d)

	client := m.(*octopusdeploy.Client)
	createdCertificate, err := client.Certificates.Add(certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenCertificate(ctx, d, createdCertificate)
	return nil
}

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Certificates.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	certificate, err := client.Certificates.GetByID(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	flattenCertificate(ctx, d, certificate)
	return nil
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate := expandCertificate(d)

	client := m.(*octopusdeploy.Client)
	updatedCertificate, err := client.Certificates.Update(*certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenCertificate(ctx, d, updatedCertificate)
	return nil
}
