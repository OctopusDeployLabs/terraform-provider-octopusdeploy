package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		DeleteContext: resourceCertificateDelete,
		Description:   "This resource manages certificates in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceCertificateRead,
		Schema:        getCertificateSchema(),
		UpdateContext: resourceCertificateUpdate,
	}
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate := expandCertificate(d)

	log.Printf("[INFO] creating certificate: %#v", certificate)

	client := m.(*octopusdeploy.Client)
	createdCertificate, err := client.Certificates.Add(certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCertificate(ctx, d, createdCertificate); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdCertificate.GetID())

	log.Printf("[INFO] certificate created (%s)", d.Id())
	return nil
}

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting certificate (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Certificates.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] certificate deleted")
	return nil
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading certificate (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	certificate, err := client.Certificates.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] certificate (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setCertificate(ctx, d, certificate); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] certificate read (%s)", d.Id())
	return nil
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating certificate (%s)", d.Id())

	certificate := expandCertificate(d)
	client := m.(*octopusdeploy.Client)
	updatedCertificate, err := client.Certificates.Update(*certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCertificate(ctx, d, updatedCertificate); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] certificate updated (%s)", d.Id())
	return nil
}
