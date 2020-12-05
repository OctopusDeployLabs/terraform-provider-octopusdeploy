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
		Description: "Provides information about existing certificates.",
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
		flattenedCertificates = append(flattenedCertificates, flattenCertificate(certificate))
	}

	d.Set("certificate", flattenedCertificates)
	d.SetId("Certificates " + time.Now().UTC().String())

	return nil
}
