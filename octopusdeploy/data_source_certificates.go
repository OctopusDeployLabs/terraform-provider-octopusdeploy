package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificates() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing certificates.",
		ReadContext: dataSourceCertificatesRead,
		Schema:      getCertificateDataSchema(),
	}
}

func dataSourceCertificatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := certificates.CertificatesQuery{
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

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	existingCertificates, err := certificates.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedCertificates := []interface{}{}
	for _, certificate := range existingCertificates.Items {
		flattenedCertificates = append(flattenedCertificates, flattenCertificate(certificate))
	}

	d.Set("certificates", flattenedCertificates)
	d.SetId("Certificates " + time.Now().UTC().String())

	return nil
}
