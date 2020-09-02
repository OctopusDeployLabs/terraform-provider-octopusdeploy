package model

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataMachinePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataMachinePolicyReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isdefault": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataMachinePolicyReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	policyName := d.Get("name").(string)
	policies, err := apiClient.MachinePolicies.GetAll()
	if err == client.ErrItemNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading machine policy with name %s: %s", policyName, err.Error())
	}

	for _, p := range *policies {
		if p.Name == policyName {
			d.SetId(p.ID)
			d.Set("description", p.Description)
			d.Set("isdefault", p.IsDefault)
		}
	}

	return nil
}
