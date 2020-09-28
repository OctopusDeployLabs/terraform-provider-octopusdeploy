package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataAwsAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataAwsAccountReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constAccountType: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  constAccountTypeAWS,
			},
			constEnvironments: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			constTenantTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
			constSecretKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constAccessKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataAwsAccountReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingAWSAccount, name, err)
	}
	if resource == nil {
		// d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)

	return nil
}
