package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataCertificateReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constNotes: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constCertificateData: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			constPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"Certificate_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
			constTenantIDs: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataCertificateReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	name := d.Get(constName).(string)
	resourceList, err := apiClient.Certificates.GetByPartialName(name)

	if err != nil {
		return createResourceOperationError(errorReadingCertificate, name, err)
	}
	if len(resourceList) == 0 {
		// d.SetId(constEmptyString)
		return nil
	}

	// NOTE: two or more certificates could have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constCertificate, m)

			d.SetId(resource.ID)
			d.Set(constName, resource.Name)

			return nil
		}
	}

	return nil
}
