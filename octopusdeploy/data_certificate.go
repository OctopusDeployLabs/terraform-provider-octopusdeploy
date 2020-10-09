package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataCertificateReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"certificate_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"tenant_ids": {
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
			"thumbprint": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataCertificateReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	CertificateName := d.Get("name")
	cert, err := client.Certificate.GetByName(CertificateName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Certificate with name %s: %s", CertificateName, err.Error())
	}

	d.SetId(cert.ID)

	d.Set("name", cert.Name)
	d.Set("serial_number", cert.SerialNumber)
	d.Set("thumbprint", cert.Thumbprint)

	return nil
}
