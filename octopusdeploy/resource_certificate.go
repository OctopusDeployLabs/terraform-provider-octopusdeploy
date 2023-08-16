package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	client := m.(*client.Client)
	createdCertificate, err := certificates.Add(client, certificate)
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

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	if err := certificates.DeleteByID(client, spaceID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] certificate deleted")
	return nil
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading certificate (%s)", d.Id())

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	certificate, err := certificates.GetByID(client, spaceID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "certificate")
	}

	if err := setCertificate(ctx, d, certificate); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] certificate read (%s)", d.Id())
	return nil
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating certificate (%s)", d.Id())

	client := m.(*client.Client)
	certificate := expandCertificate(d)
	if certificate.CertificateData.NewValue != nil {
		newCert := &certificates.ReplacementCertificate{
			CertificateData: *certificate.CertificateData.NewValue,
			Password:        *certificate.Password.NewValue,
		}
		replaceCertificate, err := client.Certificates.Replace(certificate.ID, newCert)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := setCertificate(ctx, d, replaceCertificate); err != nil {
			return diag.FromErr(err)
		}
	}
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
