package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataAccountReadByName,
		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataAccountReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingAccount, name, err)
	}
	if resource == nil {
		return nil
	}

	logResource(constAccount, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)

	return nil
}
