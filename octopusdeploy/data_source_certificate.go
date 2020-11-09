package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"notes": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"certificate_data": {
				Type:      schema.TypeString,
				Required:  true,
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
		},
	}
}

func dataSourceCertificateReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	resourceList, err := client.Certificates.GetByPartialName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resourceList) == 0 {
		return nil
	}

	// NOTE: two or more certificates could have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constCertificate, m)

			d.SetId(resource.GetID())
			d.Set("name", resource.Name)

			return nil
		}
	}

	return nil
}
