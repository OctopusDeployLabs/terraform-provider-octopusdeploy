package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataMachinePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataMachinePolicyReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
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

func dataMachinePolicyReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	resourceList, err := client.MachinePolicies.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(resourceList) == 0 {
		return nil
	}

	// NOTE: two or more machine policies could have the same name in Octopus
	// and therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			d.SetId(resource.GetID())
			d.Set("description", resource.Description)
			d.Set("isdefault", resource.IsDefault)

			return nil
		}
	}

	return nil
}
