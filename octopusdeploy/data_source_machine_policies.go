package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMachinePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMachinePoliciesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceMachinePoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	machinePolicies, err := client.MachinePolicies.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(machinePolicies) == 0 {
		return nil
	}

	// NOTE: two or more machine policies could have the same name in Octopus
	// and therefore, a better search criteria needs to be implemented below

	for _, machinePolicy := range machinePolicies {
		if machinePolicy.Name == name {
			d.Set("description", machinePolicy.Description)
			d.Set("is_default", machinePolicy.IsDefault)
			d.SetId(machinePolicy.GetID())

			return nil
		}
	}

	return nil
}
